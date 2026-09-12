const orderDetail =
    document.getElementById(
        "order-detail"
    );


function getOrderID() {

    const path =
        window.location.pathname;

    const parts =
        path.split("/");

    return parts[
        parts.length - 1
    ];
}


async function loadOrder() {

    const orderID =
        getOrderID();

    try {

        const response =
            await fetch(
                `/api/orders/${orderID}`
            );

        if (
            response.status === 401
        ) {

            window.location.href =
                "/login";

            return;
        }

        if (
            response.status === 404
        ) {

            throw new Error(
                "Pesanan tidak ditemukan."
            );
        }

        if (!response.ok) {

            throw new Error(
                "Gagal memuat pesanan."
            );
        }

        const order =
            await response.json();

        renderOrder(
            order
        );

    } catch (error) {

        orderDetail.innerHTML = `

            <p>
                ${error.message}
            </p>

        `;

    }
}


function renderOrder(
    order,
) {

    const date =
        new Date(
            order.created_at
        );

    const paymentBadge = order.payment_method
        ? `<span class="order-payment-method">${order.payment_method}</span>`
        : "";

    const paymentActions = order.status !== "paid" && order.payment_url
        ? `
            <div class="order-payment-actions">
                <a href="${order.payment_url}" target="_blank" rel="noopener noreferrer" class="button">
                    Bayar Sekarang
                </a>
            </div>
        `
        : "";

    const invoiceButton = `
        <button type="button" id="view-invoice-button" class="secondary-button">
            Lihat Struk / Invoice
        </button>
    `;

    let itemsHTML = "";

    order.items.forEach(
        (item) => {

            const readButton = order.status === "paid" && item.file_path
                ? `
                    <a href="${item.file_path}" target="_blank" rel="noopener noreferrer" class="button secondary-button">
                        Baca Ebook
                    </a>
                `
                : "";

            itemsHTML += `

                <div class="order-item">

                    <div>

                        <h2>
                            ${item.title}
                        </h2>

                        <p>
                            ${item.author}
                        </p>

                        <p>
                            Rp${item.price.toLocaleString("id-ID")}
                            ×
                            ${item.quantity}
                        </p>

                    </div>

                    <div class="order-item-actions">
                        <strong>
                            Rp${item.subtotal.toLocaleString("id-ID")}
                        </strong>
                        ${readButton}
                    </div>

                </div>

            `;

        },
    );

    orderDetail.innerHTML = `

        <div class="order-detail-header">

            <p class="eyebrow">
                ORDER #${order.id}
            </p>

            <h1>
                Pesanan berhasil dibuat 🎉
            </h1>

            <p>
                ${date.toLocaleDateString(
                    "id-ID",
                    {
                        dateStyle:
                            "full",
                    },
                )}
            </p>

        </div>


        <div class="order-items">

            ${itemsHTML}

        </div>


        <div class="order-total">

            <span>
                Total Pesanan
            </span>

            <strong>

                Rp${order.total_amount.toLocaleString("id-ID")}

            </strong>

        </div>


        <div class="order-status-detail">

            Status:
            <strong>
                ${order.status}
            </strong>
            ${paymentBadge}

            ${order.transaction_id ? `
                <p>
                    ID Transaksi: <strong>${order.transaction_id}</strong>
                </p>
            ` : ""}

            ${order.payment_method ? `
                <p>
                    Metode Pembayaran: <strong>${order.payment_method}</strong>
                </p>
            ` : ""}

        </div>

        ${paymentActions}

        <div class="order-payment-actions">
            ${invoiceButton}
        </div>

        <a
            href="/orders"
            class="button"
        >
            Lihat Semua Pesanan
        </a>

    `;

    const viewInvoiceButton = document.getElementById("view-invoice-button");

    if (viewInvoiceButton) {
        viewInvoiceButton.addEventListener("click", async () => {
            try {
                const response = await fetch(`/api/orders/${order.id}/invoice`);

                if (!response.ok) {
                    throw new Error("Gagal memuat struk atau invoice.");
                }

                const invoice = await response.json();
                openInvoiceWindow(invoice);
            } catch (error) {
                alert(error.message);
            }
        });
    }
}

function openInvoiceWindow(invoice) {
    const printWindow = window.open("", "_blank", "width=900,height=700");

    if (!printWindow) {
        alert("Popup diblokir. Izinkan popup untuk melihat struk atau invoice.");
        return;
    }

    const itemsHTML = (invoice.items || [])
        .map((item) => `
            <tr>
                <td>${item.title}</td>
                <td>${item.author}</td>
                <td>${item.quantity}</td>
                <td>Rp${Number(item.price).toLocaleString("id-ID")}</td>
                <td>Rp${Number(item.subtotal).toLocaleString("id-ID")}</td>
            </tr>
        `)
        .join("");

    printWindow.document.write(`
        <!DOCTYPE html>
        <html lang="id">
        <head>
            <meta charset="UTF-8">
            <title>Struk #${invoice.id}</title>
            <style>
                body { font-family: Arial, sans-serif; padding: 32px; color: #1c1c1c; }
                h1 { margin-bottom: 8px; }
                .meta { margin-bottom: 24px; color: #555; }
                table { width: 100%; border-collapse: collapse; margin-top: 16px; }
                th, td { border: 1px solid #ddd; padding: 10px; text-align: left; }
                th { background: #f5f1e8; }
                .total { margin-top: 24px; font-size: 20px; font-weight: bold; text-align: right; }
                @media print { body { padding: 0; } }
            </style>
        </head>
        <body>
            <h1>Struk / Invoice Bukuin</h1>
            <div class="meta">
                <p><strong>Order #:</strong> ${invoice.id}</p>
                <p><strong>Tanggal:</strong> ${new Date(invoice.created_at).toLocaleDateString("id-ID", { dateStyle: "long" })}</p>
                <p><strong>Status:</strong> ${invoice.status}</p>
                ${invoice.payment_method ? `<p><strong>Metode Pembayaran:</strong> ${invoice.payment_method}</p>` : ""}
                ${invoice.transaction_id ? `<p><strong>ID Transaksi:</strong> ${invoice.transaction_id}</p>` : ""}
            </div>
            <table>
                <thead>
                    <tr>
                        <th>Judul</th>
                        <th>Penulis</th>
                        <th>Qty</th>
                        <th>Harga</th>
                        <th>Subtotal</th>
                    </tr>
                </thead>
                <tbody>${itemsHTML}</tbody>
            </table>
            <div class="total">Total: Rp${Number(invoice.total_amount).toLocaleString("id-ID")}</div>
        </body>
        </html>
    `);

    printWindow.document.close();
    printWindow.focus();
    printWindow.print();
}


loadOrder();