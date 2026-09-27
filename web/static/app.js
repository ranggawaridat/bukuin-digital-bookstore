document.addEventListener("DOMContentLoaded", async () => {
    const navs = document.querySelectorAll(".navbar nav");

    if (!navs.length) {
        return;
    }

    let user = null;

    try {
        const response = await fetch("/api/auth/me");

        if (response.status === 401) {
            return;
        }

        if (response.ok) {
            user = await response.json();
        }
    } catch (error) {
        console.error("Gagal mengambil data pengguna:", error);
    }

    if (user && window.location.pathname === "/") {
        window.location.href = "/books";
        return;
    }

    navs.forEach((nav) => {
        if (user) {
            nav.innerHTML = `
                <a href="/cart">Keranjang</a>
                <a href="/library">Koleksi</a>
                ${user.role === "admin" ? `
                    <a href="/admin">Admin</a>
                ` : ""}
                <a href="/profile">Profil</a>
            `;
        } else {
            nav.innerHTML = `
                <a href="/login">Masuk</a>
                <a href="/register">Daftar</a>
            `;
        }

    });
});
