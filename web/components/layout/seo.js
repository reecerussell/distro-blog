import React from "react";
import PropTypes from "prop-types";
import Head from "next/head";

const Seo = ({ title, description }) => (
	<Head>
		<title>{title}</title>
		<meta name="description" content={description} />
	</Head>
);

export default Seo;
