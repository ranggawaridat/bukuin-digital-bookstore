const loginForm =
    document.getElementById("login-form");

const message =
    document.getElementById("message");


loginForm.addEventListener(
    "submit",
    async (event) => {
        event.preventDefault();

        const email =
            document.getElementById("email").value;

        const password =
            document.getElementById("password").value;

        try {

            const response = await fetch(
                "/api/auth/login",
                {
                    method: "POST",

                    headers: {
                        "Content-Type":
                            "application/json",
                    },

                    body: JSON.stringify({
                        email,
                        password,
                    }),
                },
            );

            if (!response.ok) {
                const error =
                    await response.text();

                throw new Error(error);
            }

            message.textContent =
                "Login berhasil...";

            window.location.href =
                "/books";

        } catch (error) {

            message.textContent =
                error.message;

        }

    },
);