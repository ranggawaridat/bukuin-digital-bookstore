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
                ${user.role === "admin" ? '<a href="/admin">Admin</a>' : ""}
                <details class="profile-menu">
                    <summary class="profile-trigger">Profile</summary>
                    <div class="profile-menu-panel">
                        <a href="/profile">Profile</a>
                        <a href="/orders">Orders</a>
                        <button type="button" class="logout-button" data-action="logout">Logout</button>
                    </div>
                </details>
            `;
        } else {
            nav.innerHTML = `
                <a href="/login">Login</a>
                <a href="/register">Register</a>
            `;
        }

        const logoutButton = nav.querySelector('[data-action="logout"]');

        if (logoutButton) {
            logoutButton.addEventListener("click", async () => {
                try {
                    const response = await fetch("/api/auth/logout", {
                        method: "POST",
                    });

                    if (!response.ok) {
                        throw new Error("Gagal logout.");
                    }

                    window.location.href = "/login";
                } catch (error) {
                    alert(error.message);
                }
            });
        }
    });
});
