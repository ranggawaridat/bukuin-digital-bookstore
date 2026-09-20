const cartItems =
    document.getElementById("cart-items");

const cartSummary =
    document.getElementById("cart-summary");


async function loadCart() {

    try {

        const response =
            await fetch(
                "/api/cart"
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
                "Gagal memuat keranjang."
            );
        }

        const cart =
            await response.json();

        renderCart(cart);

    } catch (error) {

        cartItems.innerHTML = `
            <p>
                ${error.message}
            </p>
        `;

        console.error(error);

    }
}


function renderCart(cart) {

    cartItems.innerHTML = "";

    if (
        !cart.items ||
        cart.items.length === 0
    ) {

        cartItems.innerHTML = `

            <div class="empty-cart">

                <h2>
                    Keranjangmu masih kosong.
                </h2>

                <p>
                    Sepertinya belum ada buku
                    yang kamu bawa pulang 📚
                </p>

                <a
                    href="/books"
                    class="button"
                >
                    Lihat Buku
                </a>

            </div>

        `;

        cartSummary.innerHTML = "";

        return;
    }

    let total = 0;

    cart.items.forEach(
        (item) => {

            const subtotal = item.price;

            total += subtotal;

            const element =
                document.createElement(
                    "article"
                );

            element.className =
                "cart-item";

            const coverMarkup = item.cover_url
                ? `<img src="${item.cover_url}" alt="${item.title}">`
                : "📖";

            element.innerHTML = `

                <div class="cart-item-cover">
                    ${coverMarkup}
                </div>

                <div class="cart-item-info">

                    <h2>
                        ${item.title}
                    </h2>

                    <p>
                        ${item.author}
                    </p>

                    <strong>
                        Rp${item.price.toLocaleString("id-ID")}
                    </strong>

                </div>

                <div class="cart-item-actions">
                    <p>
                        Rp${subtotal.toLocaleString("id-ID")}
                    </p>

                    <button
                        class="delete-button"
                        data-action="delete"
                    >
                        Hapus
                    </button>

                </div>

            `;

            const deleteButton =
                element.querySelector(
                    '[data-action="delete"]'
                );

            deleteButton.addEventListener(
                "click",
                async () => {

                    await deleteItem(
                        item.id
                    );

                },
            );

            cartItems.appendChild(
                element
            );

        },
    );

    cartSummary.innerHTML = `

        <div class="cart-total">

            <span>
                Total
            </span>

            <strong>
                Rp${total.toLocaleString("id-ID")}
            </strong>

        </div>

        <button
            id="checkout-button"
            class="button"
        >
            Checkout
        </button>

    `;

    const checkoutButton =
        document.getElementById(
            "checkout-button"
        );

    checkoutButton.addEventListener(
        "click",
        checkout,
    );
}


async function deleteItem(
    itemID,
) {

    try {

        const response =
            await fetch(
                `/api/cart/items/${itemID}`,
                {
                    method:
                        "DELETE",
                },
            );

        if (!response.ok) {
            throw new Error(
                "Gagal menghapus buku."
            );
        }

        loadCart();

    } catch (error) {

        alert(
            error.message
        );

        console.error(error);

    }
}


async function checkout() {

    const checkoutButton =
        document.getElementById(
            "checkout-button"
        );

    checkoutButton.disabled =
        true;

    checkoutButton.textContent =
        "Memproses...";

    try {

        const response =
            await fetch(
                "/api/orders/checkout",
                {
                    method: "POST",
                },
            );

        if (!response.ok) {

            const message =
                await response.text();

            throw new Error(
                message ||
                "Checkout gagal."
            );
        }

        const order =
            await response.json();

        if (order.payment_url) {
            window.location.href = order.payment_url;
            return;
        }

        window.location.href =
            `/orders/${order.id}`;

    } catch (error) {

        alert(
            error.message
        );

        console.error(error);

        checkoutButton.disabled =
            false;

        checkoutButton.textContent =
            "Checkout";

    }
}


loadCart();