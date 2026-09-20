const statsGrid = document.getElementById("stats-grid");
const booksTableBody = document.getElementById("books-table-body");
const ordersList = document.getElementById("orders-list");
const bestSellingList = document.getElementById("best-selling-list");
const bookForm = document.getElementById("book-form");
const cancelEditButton = document.getElementById("cancel-edit");
const printReportButton = document.getElementById("print-report-button");

async function requireAdmin() {
    const response = await fetch("/api/auth/me");

    if (response.status === 401) {
        window.location.href = "/login";
        return null;
    }

    if (!response.ok) {
        throw new Error("Gagal memvalidasi akses admin.");
    }

    const user = await response.json();

    if (user.role !== "admin") {
        window.location.href = "/books";
        return null;
    }

    return user;
}

async function loadDashboard() {
    try {
        const user = await requireAdmin();

        if (!user) {
            return;
        }

        const [statsResponse, booksResponse, ordersResponse] = await Promise.all([
            fetch("/api/admin/stats"),
            fetch("/api/books"),
            fetch("/api/admin/orders")
        ]);

        if (!statsResponse.ok || !ordersResponse.ok) {
            throw new Error("Gagal memuat data dashboard.");
        }

        const stats = await statsResponse.json();
        const books = await booksResponse.json();
        const orders = await ordersResponse.json();

        renderStats(stats);
        renderBooksTable(books);
        renderOrdersSummary(orders);
        renderBestSelling(books, orders);
    } catch (error) {
        statsGrid.innerHTML = `
            <p>${error.message}</p>
        `;
    }
}

function renderStats(stats) {
    statsGrid.innerHTML = `
        <article class="stat-card">
            <span>Total User</span>
            <strong>${stats.total_users}</strong>
        </article>
        <article class="stat-card">
            <span>Total Buku</span>
            <strong>${stats.total_books}</strong>
        </article>
        <article class="stat-card">
            <span>Total Pesanan</span>
            <strong>${stats.total_orders}</strong>
        </article>
        <article class="stat-card">
            <span>Revenue</span>
            <strong>Rp${Number(stats.total_revenue || 0).toLocaleString("id-ID")}</strong>
        </article>
    `;
}

function renderBooksTable(books) {
    if (!books || books.length === 0) {
        booksTableBody.innerHTML = `
            <tr>
                <td colspan="4">Belum ada buku.</td>
            </tr>
        `;
        return;
    }

    booksTableBody.innerHTML = books
        .map(
            (book) => `
                <tr>
                    <td>${book.title}</td>
                    <td>${book.category}</td>
                    <td>Rp${Number(book.price).toLocaleString("id-ID")}</td>
                    <td class="table-actions">
                        <button type="button" class="secondary-button" data-action="edit" data-id="${book.id}">
                            Edit
                        </button>
                        <button type="button" class="danger-button" data-action="delete" data-id="${book.id}">
                            Hapus
                        </button>
                    </td>
                </tr>
            `
        )
        .join("");

    booksTableBody
        .querySelectorAll("button[data-action]")
        .forEach((button) => {
            button.addEventListener("click", async () => {
                const action = button.dataset.action;
                const id = Number(button.dataset.id);

                if (action === "edit") {
                    fillFormForEdit(id);
                    return;
                }

                if (action === "delete") {
                    await deleteBook(id);
                }
            });
        });
}

function renderOrdersSummary(orders) {
    if (!orders || orders.length === 0) {
        ordersList.innerHTML = "Belum ada pesanan.";
        return;
    }

    ordersList.innerHTML = orders
        .slice(0, 6)
        .map(
            (order) => `
                <article class="admin-order-item">
                    <div>
                        <strong>#${order.id}</strong>
                        <p>${order.user_name || "User"}</p>
                        <p>
                            ${order.transaction_id ? `Transaksi: ${order.transaction_id}` : "Transaksi belum dibuat"}
                        </p>
                        ${order.payment_url ? `<p><a href="${order.payment_url}" target="_blank" rel="noopener noreferrer">Buka link pembayaran</a></p>` : ""}
                    </div>
                    <div>
                        <p>Rp${Number(order.total_amount).toLocaleString("id-ID")}</p>
                        <div class="admin-order-controls">
                            <label class="order-status-control">
                                <span>Status</span>
                                <select data-order-id="${order.id}" data-status="${order.status}">
                                    <option value="pending" ${order.status === "pending" ? "selected" : ""}>Pending</option>
                                    <option value="paid" ${order.status === "paid" ? "selected" : ""}>Paid</option>
                                    <option value="cancelled" ${order.status === "cancelled" ? "selected" : ""}>Cancelled</option>
                                    <option value="refunded" ${order.status === "refunded" ? "selected" : ""}>Refunded</option>
                                </select>
                            </label>
                            <button type="button" class="secondary-button" data-invoice-order="${order.id}">Invoice</button>
                        </div>
                        ${order.payment_method ? `<span class="payment-tag">${order.payment_method}</span>` : ""}
                    </div>
                </article>
            `
        )
        .join("");

    ordersList.querySelectorAll("[data-invoice-order]").forEach((button) => {
        button.addEventListener("click", async () => {
            const orderID = Number(button.dataset.invoiceOrder);
            const response = await fetch(`/api/admin/orders/${orderID}/invoice`);
            if (!response.ok) {
                alert("Gagal memuat invoice.");
                return;
            }
            const invoice = await response.json();
            openInvoiceWindow(invoice);
        });
    });

    ordersList.querySelectorAll("select[data-order-id]").forEach((select) => {
        select.addEventListener("change", async (event) => {
            const orderID = Number(event.target.dataset.orderId);
            const status = event.target.value;

            try {
                const response = await fetch(`/api/admin/orders/${orderID}/status`, {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ status })
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(errorText || "Gagal mengubah status order.");
                }

                await loadDashboard();
            } catch (error) {
                alert(error.message);
                event.target.value = event.target.dataset.status;
            }
        });
    });
}

function openInvoiceWindow(invoice) {
    const printWindow = window.open("", "_blank", "width=900,height=700");
    if (!printWindow) {
        alert("Popup diblokir. Izinkan popup untuk melihat invoice.");
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
            <title>Invoice #${invoice.id}</title>
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
            <h1>Invoice Bukuin</h1>
            <div class="meta">
                <p><strong>Invoice #:</strong> ${invoice.id}</p>
                <p><strong>Tanggal:</strong> ${new Date(invoice.created_at).toLocaleDateString("id-ID", { dateStyle: "long" })}</p>
                <p><strong>Pelanggan:</strong> ${invoice.user_name} (${invoice.user_email || "-"})</p>
                <p><strong>Status:</strong> ${invoice.status}</p>
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

function openReportWindow(report) {
    const printWindow = window.open("", "_blank", "width=1100,height=800");
    if (!printWindow) {
        alert("Popup diblokir. Izinkan popup untuk melihat laporan.");
        return;
    }

    const stats = report.stats || {};
    const orders = report.orders || [];

    const orderRows = orders.length
        ? orders
              .map(
                  (order) => `
                    <tr>
                        <td>#${order.id}</td>
                        <td>${order.user_name || "-"}</td>
                        <td>${order.user_email || "-"}</td>
                        <td>Rp${Number(order.total_amount).toLocaleString("id-ID")}</td>
                        <td>${order.status}</td>
                        <td>${new Date(order.created_at).toLocaleDateString("id-ID", { dateStyle: "medium" })}</td>
                    </tr>
                `
              )
              .join("")
        : `
            <tr>
                <td colspan="6">Belum ada data pesanan.</td>
            </tr>
        `;

    printWindow.document.write(`
        <!DOCTYPE html>
        <html lang="id">
        <head>
            <meta charset="UTF-8">
            <title>Laporan Bukuin</title>
            <style>
                body { font-family: Arial, sans-serif; padding: 32px; color: #1c1c1c; }
                h1 { margin-bottom: 8px; }
                .meta { margin-bottom: 24px; color: #555; }
                .stats { display: grid; grid-template-columns: repeat(4, minmax(180px, 1fr)); gap: 12px; margin: 20px 0; }
                .stat { border: 1px solid #ddd; padding: 12px; background: #f9f9f9; }
                .stat strong { display: block; font-size: 24px; margin-top: 8px; }
                table { width: 100%; border-collapse: collapse; margin-top: 16px; }
                th, td { border: 1px solid #ddd; padding: 10px; text-align: left; }
                th { background: #f5f1e8; }
                @media print { body { padding: 0; } }
            </style>
        </head>
        <body>
            <h1>Laporan Bukuin</h1>
            <div class="meta">
                <p><strong>Generated:</strong> ${new Date(report.generated_at).toLocaleString("id-ID")}</p>
            </div>
            <div class="stats">
                <div class="stat">
                    <span>Total User</span>
                    <strong>${stats.total_users || 0}</strong>
                </div>
                <div class="stat">
                    <span>Total Buku</span>
                    <strong>${stats.total_books || 0}</strong>
                </div>
                <div class="stat">
                    <span>Total Pesanan</span>
                    <strong>${stats.total_orders || 0}</strong>
                </div>
                <div class="stat">
                    <span>Revenue</span>
                    <strong>Rp${Number(stats.total_revenue || 0).toLocaleString("id-ID")}</strong>
                </div>
            </div>
            <table>
                <thead>
                    <tr>
                        <th>#Order</th>
                        <th>Pelanggan</th>
                        <th>Email</th>
                        <th>Total</th>
                        <th>Status</th>
                        <th>Tanggal</th>
                    </tr>
                </thead>
                <tbody>${orderRows}</tbody>
            </table>
        </body>
        </html>
    `);
    printWindow.document.close();
    printWindow.focus();
    printWindow.print();
}

function renderBestSelling(books, orders) {
    if (!books || books.length === 0) {
        bestSellingList.innerHTML = "Belum ada data buku.";
        return;
    }

    const salesByBook = new Map();

    orders.forEach((order) => {
        if (!order.items) {
            return;
        }

        order.items.forEach((item) => {
            salesByBook.set(
                item.book_id,
                (salesByBook.get(item.book_id) || 0) + item.quantity
            );
        });
    });

    const sortedBooks = [...books]
        .sort((a, b) => (salesByBook.get(b.id) || 0) - (salesByBook.get(a.id) || 0))
        .slice(0, 5);

    bestSellingList.innerHTML = sortedBooks
        .map((book, index) => `
            <li>
                <span>${index + 1}. ${book.title}</span>
                <strong>${salesByBook.get(book.id) || 0} terjual</strong>
            </li>
        `)
        .join("");
}

async function fillFormForEdit(id) {
    const response = await fetch(`/api/books/${id}`);

    if (!response.ok) {
        alert("Gagal memuat data buku.");
        return;
    }

    const book = await response.json();

    document.getElementById("book-id").value = book.id;
    document.getElementById("book-title").value = book.title;
    document.getElementById("book-author").value = book.author;
    document.getElementById("book-category").value = book.category;
    document.getElementById("book-price").value = book.price;
    document.getElementById("book-description").value = book.description;

    window.scrollTo({ top: 0, behavior: "smooth" });
}

async function deleteBook(id) {
    const confirmed = window.confirm("Apakah Anda yakin ingin menghapus buku ini?");

    if (!confirmed) {
        return;
    }

    const response = await fetch(`/api/admin/books/${id}`, {
        method: "DELETE"
    });

    if (!response.ok) {
        alert("Gagal menghapus buku.");
        return;
    }

    await loadDashboard();
}

bookForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const formData = new FormData(bookForm);
    const bookID = document.getElementById("book-id").value;
    const method = bookID ? "PUT" : "POST";
    const url = bookID ? `/api/admin/books/${bookID}` : "/api/admin/books";

    const body = new FormData();

    body.append("title", String(formData.get("title") || ""));
    body.append("author", String(formData.get("author") || ""));
    body.append("category", String(formData.get("category") || ""));
    body.append("description", String(formData.get("description") || ""));
    body.append("price", String(formData.get("price") || "0"));

    const coverFile = formData.get("cover");
    if (coverFile instanceof File && coverFile.size > 0) {
        body.append("cover", coverFile);
    }

    const ebookFile = formData.get("ebook");
    if (ebookFile instanceof File && ebookFile.size > 0) {
        body.append("ebook", ebookFile);
    }

    const response = await fetch(url, {
        method,
        body
    });

    if (!response.ok) {
        const errorText = await response.text();
        alert(errorText || "Gagal menyimpan buku.");
        return;
    }

    bookForm.reset();
    document.getElementById("book-id").value = "";
    await loadDashboard();
});

cancelEditButton.addEventListener("click", () => {
    bookForm.reset();
    document.getElementById("book-id").value = "";
});

if (printReportButton) {
    printReportButton.addEventListener("click", async () => {
        try {
            const response = await fetch("/api/admin/report");

            if (!response.ok) {
                throw new Error("Gagal memuat laporan.");
            }

            const report = await response.json();
            openReportWindow(report);
        } catch (error) {
            alert(error.message);
        }
    });
}

loadDashboard();
