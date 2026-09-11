const ordersList =
    document.getElementById(
        "orders-list"
    );


async function loadOrders() {

    try {

        const response =
            await fetch(
                "/api/orders"
            );

        if (
            response.status === 401
        ) {

            window.location.href =
                "/login";

            return;
        }

        if (!response.ok) {

            throw new Error(
                "Gagal memuat pesanan."
            );
        }

        const orders =
            await response.json();

        renderOrders(
            orders
        );

    } catch (error) {

        ordersList.innerHTML = `
            <p>
                ${error.message}
            </p>
        `;

    }
}


function renderOrders(
    orders,
) {

    ordersList.innerHTML = "";

    if (
        orders.length === 0
    ) {

        ordersList.innerHTML = `

            <div class="empty-cart">

                <h2>
                    Belum ada pesanan.
                </h2>

                <p>
                    Yuk cari buku pertama kamu 📚
                </p>

                <a
                    href="/books"
                    class="button"
                >
                    Lihat Buku
                </a>

            </div>

        `;

        return;
    }

    orders.forEach(
        (order) => {

            const element =
                document.createElement(
                    "article"
                );

            element.className =
                "order-card";

            const date =
                new Date(
                    order.created_at
                );

            element.innerHTML = `

                <div>

                    <p class="order-number">
                        ORDER #${order.id}
                    </p>

                    <h2>
                        Rp${order.total_amount.toLocaleString("id-ID")}
                    </h2>

                    <p>
                        ${date.toLocaleDateString(
                            "id-ID",
                            {
                                dateStyle:
                                    "long",
                            },
                        )}
                    </p>

                </div>

                <div class="order-status">

                    ${order.status}

                </div>

            `;

            element.addEventListener(
                "click",
                () => {

                    window.location.href =
                        `/orders/${order.id}`;

                },
            );

            ordersList.appendChild(
                element
            );

        },
    );
}


loadOrders();