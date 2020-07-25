import React, { useState, useEffect } from "react";
import Layout from "../components/layout";
import InitClient from "../components/client";
import { parse as parseUrl } from "url";
import fetch from "isomorphic-unfetch";
import Error from "next/error";

const Index = ({ status, data }) => {
	if (status !== 200) {
		return <Error statusCode={status} />;
	}

	const [isNavigationOpen, setNavigationStatus] = useState(false);
	const toggleNavigation = () => setNavigationStatus(!isNavigationOpen);

	useEffect(() => InitClient(), []);

	return (
		<Layout
			isNavigationOpen={isNavigationOpen}
			toggleNavigation={toggleNavigation}
			data={data}
		>
			<div dangerouslySetInnerHTML={{ __html: data.content }}></div>
		</Layout>
	);
};

const getServerSideProps = async ({ req }) => {
	const { pathname } = parseUrl(req.url, true);
	const res = await fetch(
		"https://api.reece-russell.co.uk/ui/page/" + pathname.substring(1),
		{
			headers: {
				"x-api-key": "xYC4t7bmkO2I3fdOIlRzAWkbpIBVSMnakIWF1vl1",
			},
		}
	);
	let resData = null;

	try {
		resData = await res.json();
	} catch (err) {
		console.log(err);
		return {
			props: { status: 500 },
		};
	}

	if (res.status !== 200 && resData.error) {
		console.error(resData.error);
	}

	return {
		props: {
			status: res.status,
			...resData,
		},
	};
};

export default Index;
export { getServerSideProps };
