const profileName = document.getElementById("profile-name");
const profileEmail = document.getElementById("profile-email");
const profileRole = document.getElementById("profile-role");
const profileCreatedAt = document.getElementById("profile-created-at");
const profileOrdersList = document.getElementById("profile-orders-list");
const logoutButton = document.getElementById("logout-button");

async function loadProfile() {
    try {
        const meResponse = await fetch("/api/auth/me");

        if (meResponse.status === 401) {
            window.location.href = "/login";
            return;
        }

        if (!meResponse.ok) {
            throw new Error("Gagal memuat profil.");
        }

        const user = await meResponse.json();

        profileName.textContent = user.name;
        profileEmail.textContent = user.email;
        profileRole.textContent = user.role;
        profileCreatedAt.textContent = new Date(user.created_at).toLocaleDateString("id-ID", {
            dateStyle: "long",
        });

        const ordersResponse = await fetch("/api/orders");

        if (!ordersResponse.ok) {
            throw new Error("Gagal memuat riwayat pesanan.");
        }

        const orders = await ordersResponse.json();

        if (!orders || orders.length === 0) {
            profileOrdersList.innerHTML = "<p>Belum ada pesanan.</p>";
            return;
        }

        profileOrdersList.innerHTML = orders
            .map(
                (order) => `
                    <article class="profile-order-item">
                        <div>
                            <p class="eyebrow">ORDER #${order.id}</p>
                            <strong>Rp${Number(order.total_amount).toLocaleString("id-ID")}</strong>
                        </div>
                        <span class="order-status">${order.status}</span>
                    </article>
                `
            )
            .join("");
    } catch (error) {
        profileName.textContent = "Gagal memuat profil";
        profileOrdersList.innerHTML = `<p>${error.message}</p>`;
        console.error(error);
    }
}

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

loadProfile();
