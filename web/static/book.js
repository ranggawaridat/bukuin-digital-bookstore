const bookDetail =
    document.getElementById("book-detail");

const loading =
    document.getElementById("loading");


function getBookID() {
    const path =
        window.location.pathname;

    const parts =
        path.split("/");

    return parts[parts.length - 1];
}


async function loadBook() {
    const bookID =
        getBookID();

    try {

        const response =
            await fetch(
                `/api/books/${bookID}`
            );

        if (!response.ok) {

            if (
                response.status === 404
            ) {
                throw new Error(
                    "Buku tidak ditemukan."
                );
            }

            throw new Error(
                "Gagal memuat buku."
            );
        }

        const book =
            await response.json();

        renderBook(book);

    } catch (error) {

        loading.textContent =
            error.message;

        console.error(error);

    }
}


function renderBook(book) {

    loading.remove();

    document.title =
        `${book.title} - Bukuin`;

    bookDetail.innerHTML = `

        <div class="book-detail-cover">
            📖
        </div>

        <div class="book-detail-content">

            <p class="book-category">
                ${book.category}
            </p>

            <h1>
                ${book.title}
            </h1>

            <p class="book-detail-author">
                oleh ${book.author}
            </p>

            <p class="book-detail-description">
                ${book.description}
            </p>

            <p class="book-detail-price">
                Rp${book.price.toLocaleString("id-ID")}
            </p>

            <button
                id="buy-button"
                class="button"
            >
                Beli Sekarang
            </button>

        </div>

    `;

    const buyButton =
        document.getElementById(
            "buy-button"
        );

    buyButton.addEventListener(
        "click",
        () => {
            alert(
                "Shopping cart akan hadir pada Phase 6 😈"
            );
        },
    );
}


loadBook();