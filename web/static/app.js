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

    navs.forEach((nav) => {
        if (user) {
            nav.innerHTML = `
                <a href="/cart" class="nav-icon-link" aria-label="Cart" title="Cart">
                    <span aria-hidden="true">🛒</span>
                </a>
                <a href="/library" class="nav-icon-link" aria-label="My Library" title="My Library">
                    <span aria-hidden="true">📚</span>
                </a>
                ${user.role === "admin" ? `
                    <a href="/admin" class="nav-icon-link" aria-label="Admin" title="Admin">
                        <span aria-hidden="true">⚙️</span>
                    </a>
                ` : ""}
                <a href="/profile" class="nav-icon-link" aria-label="Profile" title="Profile">
                    <span aria-hidden="true">👤</span>
                </a>
            `;
        } else {
            nav.innerHTML = `
                <a href="/login">Login</a>
                <a href="/register">Register</a>
            `;
        }

    });
});
