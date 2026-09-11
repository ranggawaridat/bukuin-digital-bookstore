const orderDetail =
    document.getElementById(
        "order-detail"
    );


function getOrderID() {

    const path =
        window.location.pathname;

    const parts =
        path.split("/");

    return parts[
        parts.length - 1
    ];
}


async function loadOrder() {

    const orderID =
        getOrderID();

    try {

        const response =
            await fetch(
                `/api/orders/${orderID}`
            );

        if (
            response.status === 401
        ) {

            window.location.href =
                "/login";

            return;
        }

        if (
            response.status === 404
        ) {

            throw new Error(
                "Pesanan tidak ditemukan."
            );
        }

        if (!response.ok) {

            throw new Error(
                "Gagal memuat pesanan."
            );
        }

        const order =
            await response.json();

        renderOrder(
            order
        );

    } catch (error) {

        orderDetail.innerHTML = `

            <p>
                ${error.message}
            </p>

        `;

    }
}


function renderOrder(
    order,
) {

    const date =
        new Date(
            order.created_at
        );

    let itemsHTML = "";

    order.items.forEach(
        (item) => {

            itemsHTML += `

                <div class="order-item">

                    <div>

                        <h2>
                            ${item.title}
                        </h2>

                        <p>
                            ${item.author}
                        </p>

                        <p>
                            Rp${item.price.toLocaleString("id-ID")}
                            ×
                            ${item.quantity}
                        </p>

                    </div>

                    <strong>

                        Rp${item.subtotal.toLocaleString("id-ID")}

                    </strong>

                </div>

            `;

        },
    );

    orderDetail.innerHTML = `

        <div class="order-detail-header">

            <p class="eyebrow">
                ORDER #${order.id}
            </p>

            <h1>
                Pesanan berhasil dibuat 🎉
            </h1>

            <p>
                ${date.toLocaleDateString(
                    "id-ID",
                    {
                        dateStyle:
                            "full",
                    },
                )}
            </p>

        </div>


        <div class="order-items">

            ${itemsHTML}

        </div>


        <div class="order-total">

            <span>
                Total Pesanan
            </span>

            <strong>

                Rp${order.total_amount.toLocaleString("id-ID")}

            </strong>

        </div>


        <div class="order-status-detail">

            Status:
            <strong>
                ${order.status}
            </strong>

        </div>


        <a
            href="/orders"
            class="button"
        >
            Lihat Semua Pesanan
        </a>

    `;
}


loadOrder();