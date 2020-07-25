import React from "react";
import { Container, Row, Col } from "reactstrap";

const Header = ({ title, description, imageId }) => {
	let style = {};
	if (imageId) {
		style.backgroundImage = `url(https://api.reece-russell.co.uk/media/${imageId})`;
	}

	return (
		<header className="masthead" style={style}>
			<div className="overlay"></div>
			<Container>
				<Row>
					<Col lg="8" md="10" className="mx-auto">
						<div className="site-heading">
							<h1>{title}</h1>
							<span className="subheading">{description}</span>
						</div>
					</Col>
				</Row>
			</Container>
		</header>
	);
};

export default Header;
