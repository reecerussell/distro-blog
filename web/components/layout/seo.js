import React from "react";
import PropTypes from "prop-types";
import Head from "next/head";

const Seo = ({
	title,
	description,
	siteName,
	index,
	follow,
	isBlog,
	imageId,
}) => {
	let robots = "";
	if (!index && !follow) {
		robots = "none";
	} else if (!index) {
		robots = "noindex";
	} else if (!follow) {
		robots = "nofollow";
	} else {
		robots = "all";
	}

	return (
		<Head>
			<title>{title}</title>
			<meta name="description" content={description} />
			<meta name="robots" content={robots} />
			<meta name="googlebot" content={robots} />

			<meta property="og:type" content={isBlog ? "blog" : "website"} />
			<meta property="og:title" content={title} />
			<meta property="og:description" content={description} />
			<meta property="og:site_name" content={siteName} />
			{imageId ? (
				<meta
					property="og:image"
					content={`https://api.reece-russell.co.uk/media/${imageId}`}
				/>
			) : null}

			<meta name="twitter:card" content="summary" />
			<meta name="twitter:title" content={title} />
			<meta name="twitter:description" content={description} />
			{imageId ? (
				<meta
					name="twitter:image"
					content={`https://api.reece-russell.co.uk/media/${imageId}`}
				/>
			) : null}
		</Head>
	);
};

Seo.propTypes = {
	title: PropTypes.string.isRequired,
	description: PropTypes.string.isRequired,
	siteName: PropTypes.string.isRequired,
	index: PropTypes.bool.isRequired,
	follow: PropTypes.bool.isRequired,
	isBlog: PropTypes.bool.isRequired,
	imageId: PropTypes.string,
};

export default Seo;
