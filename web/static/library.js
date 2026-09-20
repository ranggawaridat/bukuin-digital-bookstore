const libraryList = document.getElementById("library-list");
const libraryCount = document.getElementById("library-count");

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
            libraryCount.textContent = "Belum ada koleksi. Temukan bacaan pertamamu hari ini.";
            libraryList.innerHTML = `
                <div class="library-empty">
                    <h2>Belum ada buku di library.</h2>
                    <p>Yuk beli buku pertama kamu dan mulai baca ebooknya.</p>
                    <a href="/books" class="button">Lihat Buku</a>
                </div>
            `;
            return;
        }

        libraryCount.textContent = `${library.length} ebook siap dibaca kapan saja.`;
        libraryList.innerHTML = library
            .map((item) => {
                const paidAt = item.paid_at
                    ? new Date(item.paid_at).toLocaleDateString("id-ID", { dateStyle: "long" })
                    : "Tanggal pembelian tidak tersedia";
                const coverMarkup = item.cover_url
                    ? `<img src="${item.cover_url}" alt="Cover ${item.title}">`
                    : "📖";
                const readMarkup = item.file_path
                    ? `<a href="${item.file_path}" target="_blank" rel="noopener noreferrer" class="button">Baca Ebook</a>`
                    : `<span class="library-unavailable">File ebook belum tersedia</span>`;

                return `
                <article class="library-item">
                    <div class="library-cover">
                        ${coverMarkup}
                    </div>
                    <div class="library-content">
                        <p class="book-category">${item.category}</p>
                        <h2>${item.title}</h2>
                        <p>oleh ${item.author}</p>
                        <p class="library-paid-at">Dibeli ${paidAt}</p>
                        ${readMarkup}
                    </div>
                </article>
                `;
            })
            .join("");
    } catch (error) {
        libraryCount.textContent = "Library sedang tidak dapat diakses.";
        libraryList.innerHTML = `<p class="library-error">${error.message}</p>`;
    }
}

loadLibrary();
