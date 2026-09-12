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
                <a href="/">Home</a>
                <a href="/books">Books</a>
                <a href="/cart">Cart</a>
                <a href="/orders">Orders</a>
                ${user.role === "admin" ? '<a href="/admin">Admin</a>' : ""}
                <a href="/profile">Profile</a>
                <button type="button" class="logout-button" data-action="logout">Logout</button>
            `;
        } else {
            nav.innerHTML = `
                <a href="/">Home</a>
                <a href="/books">Books</a>
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
