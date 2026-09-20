const registerForm =
    document.getElementById("register-form");

const message =
    document.getElementById("message");


registerForm.addEventListener(
    "submit",
    async (event) => {
        event.preventDefault();

        const name =
            document.getElementById("name").value;

        const email =
            document.getElementById("email").value;

        const password =
            document.getElementById("password").value;

        try {
            const response = await fetch(
                "/api/auth/register",
                {
                    method: "POST",

                    headers: {
                        "Content-Type":
                            "application/json",
                    },

                    body: JSON.stringify({
                        name,
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
                "Akun berhasil dibuat. Mengarahkan ke login...";

            registerForm.reset();

            setTimeout(
                () => {
                    window.location.href =
                        "/login";
                },
                1000,
            );

        } catch (error) {

            message.textContent =
                error.message;

        }

    },
);