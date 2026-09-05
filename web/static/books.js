const bookList = document.getElementById("book-list");

const searchInput = document.getElementById("search");

const categorySelect = document.getElementById("category");


async function loadBooks() {
    const search = searchInput.value;

    const category = categorySelect.value;

    const params = new URLSearchParams();

    if (search) {
        params.set("q", search);
    }

    if (category) {
        params.set("category", category);
    }

    try {
        const response = await fetch(
            `/api/books?${params.toString()}`
        );

        if (!response.ok) {
            throw new Error(
                "Failed to load books"
            );
        }

        const books = await response.json();

        renderBooks(books);

    } catch (error) {
        bookList.innerHTML = `
            <p>
                Gagal memuat buku.
            </p>
        `;

        console.error(error);
    }
}


function renderBooks(books) {
    bookList.innerHTML = "";

    if (!books || books.length === 0) {

        bookList.innerHTML = `
            <p>
                Buku tidak ditemukan.
            </p>
        `;

        return;
    }

    books.forEach((book) => {

        const card =
            document.createElement("article");

        card.className = "book-card";

        card.innerHTML = `

            <div class="book-cover">
                📖
            </div>

            <p class="book-category">
                ${book.category}
            </p>

            <h2>
                ${book.title}
            </h2>

            <p class="book-author">
                ${book.author}
            </p>

            <strong class="book-price">
                Rp${book.price.toLocaleString("id-ID")}
            </strong>

        `;

        bookList.appendChild(card);

    });
}


searchInput.addEventListener(
    "input",
    () => {
        loadBooks();
    }
);


categorySelect.addEventListener(
    "change",
    () => {
        loadBooks();
    }
);


loadBooks();