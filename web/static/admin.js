const statsGrid = document.getElementById("stats-grid");
const booksTableBody = document.getElementById("books-table-body");
const ordersList = document.getElementById("orders-list");
const bestSellingList = document.getElementById("best-selling-list");
const bookForm = document.getElementById("book-form");
const cancelEditButton = document.getElementById("cancel-edit");

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
                    </div>
                    <div>
                        <p>Rp${Number(order.total_amount).toLocaleString("id-ID")}</p>
                        <span>${order.status}</span>
                    </div>
                </article>
            `
        )
        .join("");
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
    document.getElementById("book-cover-url").value = book.cover_url || "";
    document.getElementById("book-file-path").value = book.file_path || "";

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
    const payload = {
        title: formData.get("title"),
        author: formData.get("author"),
        category: formData.get("category"),
        description: formData.get("description"),
        price: Number(formData.get("price")),
        cover_url: formData.get("cover_url") || "",
        file_path: formData.get("file_path") || ""
    };

    const bookID = document.getElementById("book-id").value;
    const method = bookID ? "PUT" : "POST";
    const url = bookID ? `/api/admin/books/${bookID}` : "/api/admin/books";

    const response = await fetch(url, {
        method,
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
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

loadDashboard();
