import React, { useState, useEffect } from "react";
import Layout from "../components/layout";
import InitClient from "../components/client";

const Index = () => {
	const [isNavigationOpen, setNavigationStatus] = useState(false);
	const toggleNavigation = () => setNavigationStatus(!isNavigationOpen);

	useEffect(() => InitClient(), []);

	return (
		<Layout
			isNavigationOpen={isNavigationOpen}
			toggleNavigation={toggleNavigation}
		>
			<h1>Hello World</h1>
			<p>
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce
				vestibulum cursus elit, non malesuada magna faucibus vitae. Sed
				sed nulla eget odio gravida euismod in vel sem. Vestibulum et
				felis eu libero vehicula volutpat. Sed ultricies id nisi nec
				laoreet. Sed finibus et augue a ullamcorper. Mauris dapibus
				efficitur enim, nec ullamcorper neque. Vivamus quam erat,
				rhoncus eget ex eget, mattis accumsan elit. Maecenas leo nibh,
				pellentesque sit amet consequat eget, imperdiet quis quam.
			</p>
			<p>
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce
				vestibulum cursus elit, non malesuada magna faucibus vitae. Sed
				sed nulla eget odio gravida euismod in vel sem. Vestibulum et
				felis eu libero vehicula volutpat. Sed ultricies id nisi nec
				laoreet. Sed finibus et augue a ullamcorper. Mauris dapibus
				efficitur enim, nec ullamcorper neque. Vivamus quam erat,
				rhoncus eget ex eget, mattis accumsan elit. Maecenas leo nibh,
				pellentesque sit amet consequat eget, imperdiet quis quam.
			</p>
			<p>
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce
				vestibulum cursus elit, non malesuada magna faucibus vitae. Sed
				sed nulla eget odio gravida euismod in vel sem. Vestibulum et
				felis eu libero vehicula volutpat. Sed ultricies id nisi nec
				laoreet. Sed finibus et augue a ullamcorper. Mauris dapibus
				efficitur enim, nec ullamcorper neque. Vivamus quam erat,
				rhoncus eget ex eget, mattis accumsan elit. Maecenas leo nibh,
				pellentesque sit amet consequat eget, imperdiet quis quam.
			</p>
			<p>
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce
				vestibulum cursus elit, non malesuada magna faucibus vitae. Sed
				sed nulla eget odio gravida euismod in vel sem. Vestibulum et
				felis eu libero vehicula volutpat. Sed ultricies id nisi nec
				laoreet. Sed finibus et augue a ullamcorper. Mauris dapibus
				efficitur enim, nec ullamcorper neque. Vivamus quam erat,
				rhoncus eget ex eget, mattis accumsan elit. Maecenas leo nibh,
				pellentesque sit amet consequat eget, imperdiet quis quam.
			</p>
			<p>
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce
				vestibulum cursus elit, non malesuada magna faucibus vitae. Sed
				sed nulla eget odio gravida euismod in vel sem. Vestibulum et
				felis eu libero vehicula volutpat. Sed ultricies id nisi nec
				laoreet. Sed finibus et augue a ullamcorper. Mauris dapibus
				efficitur enim, nec ullamcorper neque. Vivamus quam erat,
				rhoncus eget ex eget, mattis accumsan elit. Maecenas leo nibh,
				pellentesque sit amet consequat eget, imperdiet quis quam.
			</p>
		</Layout>
	);
};

export default Index;
