const libraryList = document.getElementById("library-list");

async function loadLibrary() {
    try {
        const response = await fetch("/api/library");

        if (response.status === 401) {
            window.location.href = "/login";
            return;
        }

        if (!response.ok) {
            throw new Error("Gagal memuat My Library.");
        }

        const library = await response.json();

        if (!library || library.length === 0) {
            libraryList.innerHTML = `
                <div class="empty-cart">
                    <h2>Belum ada buku di library.</h2>
                    <p>Yuk beli buku pertama kamu dan mulai baca ebooknya.</p>
                    <a href="/books" class="button">Lihat Buku</a>
                </div>
            `;
            return;
        }

        libraryList.innerHTML = library
            .map((item) => `
                <article class="library-item">
                    <div class="library-cover">
                        ${item.cover_url ? `<img src="${item.cover_url}" alt="${item.title}">` : "📖"}
                    </div>
                    <div class="library-content">
                        <p class="book-category">${item.category}</p>
                        <h2>${item.title}</h2>
                        <p>oleh ${item.author}</p>
                        <p class="library-paid-at">Pembelian: ${new Date(item.paid_at).toLocaleDateString("id-ID", { dateStyle: "long" })}</p>
                        <a href="${item.file_path}" target="_blank" rel="noopener noreferrer" class="button">Baca Ebook</a>
                    </div>
                </article>
            `)
            .join("");
    } catch (error) {
        libraryList.innerHTML = `<p>${error.message}</p>`;
    }
}

loadLibrary();
