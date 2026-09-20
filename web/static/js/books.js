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
        const card = document.createElement("article");
        card.className = "book-card";

        // Saya ganti icon buku menjadi teks "BOOK" jika tidak ada gambar, 
        // agar sesuai dengan desain referensi.
        const coverMarkup = book.cover_url
            ? `<img src="${book.cover_url}" alt="${book.title}">`
            : "BOOK"; 

        // Update struktur HTML di sini
        card.innerHTML = `
            <div class="book-cover">
                ${coverMarkup}
            </div>

            <div class="book-info">
                <span class="book-category">
                    ${book.category}
                </span>

                <h2>
                    ${book.title}
                </h2>

                <p class="book-author">
                    ${book.author}
                </p>

                <div class="price-row">
                    <strong class="book-price">
                        Rp${book.price.toLocaleString("id-ID")}
                    </strong>
                    <button class="btn-beli">Beli</button>
                </div>
            </div>
        `;

        // Menambahkan event listener agar kalau di-klik pindah halaman
        card.addEventListener(
            "click",
            () => {
                window.location.href = `/books/${book.id}`;
            },
        );

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