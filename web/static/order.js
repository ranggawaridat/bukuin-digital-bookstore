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

    let itemsHTML = "";

    order.items.forEach(
        (item) => {

            const downloadButton = order.status === "paid" && item.file_path
                ? `
                    <a href="${item.file_path}" target="_blank" rel="noopener noreferrer" class="button secondary-button">
                        Download Ebook
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
                        ${downloadButton}
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

        <a
            href="/orders"
            class="button"
        >
            Lihat Semua Pesanan
        </a>

    `;
}


loadOrder();