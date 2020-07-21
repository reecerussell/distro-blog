const MQL = 992;
let previousPosition = 0;

const scrollHandler = (nav, navHeight) => {
	const currentTop =
		window.pageYOffset ||
		(document.documentElement || document.body.parentNode || document.body)
			.scrollTop;
	if (currentTop < previousPosition) {
		if (currentTop > 0 && nav.classList.contains("is-fixed")) {
			nav.classList.add("is-visible");
		} else {
			nav.classList.remove("is-visible", "is-fixed");
		}
	} else if (currentTop > previousPosition) {
		nav.classList.remove("is-fixed");

		if (currentTop > navHeight && !nav.classList.contains("is-fixed")) {
			nav.classList.add("is-fixed");
		}
	}

	previousPosition = currentTop;
};

const InitialiseNavigation = () => {
	if (window.innerWidth > MQL) {
		const mainNav = document.getElementById("mainNav");
		const navHeight = mainNav.getBoundingClientRect().height;
		window.addEventListener("scroll", () =>
			scrollHandler(mainNav, navHeight)
		);
	}
};

export default InitialiseNavigation;
