document.addEventListener("DOMContentLoaded", () => {
    const nav = document.querySelector(".navbar-links");
    const toggle = document.querySelector(".nav-toggle");
    const menu = nav ? nav.querySelector("ul") : null;

    if (!nav || !toggle) {
        return;
    }

    const syncMobileNav = () => {
        const isMobile = window.innerWidth <= 700;
        toggle.style.display = isMobile ? "inline-flex" : "none";
        toggle.style.opacity = "1";
        toggle.style.visibility = "visible";

        if (!isMobile) {
            nav.classList.remove("open");
            toggle.setAttribute("aria-expanded", "false");
            if (menu) {
                menu.style.display = "flex";
            }
            return;
        }

        const isOpen = nav.classList.contains("open");
        toggle.setAttribute("aria-expanded", String(isOpen));
        if (menu) {
            menu.style.display = isOpen ? "flex" : "none";
            menu.style.left = "50%";
            menu.style.top = "calc(100% + 12px)";
            menu.style.width = "min(260px, calc(100vw - 36px))";
            menu.style.transform = "translateX(-50%)";
            menu.style.right = "auto";
        }
    };

    const setMenuState = (isOpen) => {
        if (window.innerWidth > 700) {
            nav.classList.remove("open");
            toggle.setAttribute("aria-expanded", "false");
            if (menu) menu.style.display = "flex";
            return;
        }

        nav.classList.toggle("open", isOpen);
        toggle.setAttribute("aria-expanded", String(isOpen));

        if (menu) {
            menu.style.display = isOpen ? "flex" : "none";
            menu.style.left = "50%";
            menu.style.top = "calc(100% + 12px)";
            menu.style.width = "min(260px, calc(100vw - 36px))";
            menu.style.transform = "translateX(-50%)";
            menu.style.right = "auto";
        }
    };

    toggle.addEventListener("click", (event) => {
        event.preventDefault();
        event.stopPropagation();
        setMenuState(!nav.classList.contains("open"));
    });

    nav.querySelectorAll("a").forEach((link) => {
        link.addEventListener("click", () => {
            if (window.innerWidth <= 700) {
                setMenuState(false);
            }
        });
    });

    window.addEventListener("resize", () => {
        if (window.innerWidth > 700) {
            setMenuState(false);
        } else {
            syncMobileNav();
        }
    });

    syncMobileNav();
});
