import React from "react";
import { Container, Row, Col } from "reactstrap";

const Header = () => (
	<header
		className="masthead"
		style={{ "background-image": "url('img/home-bg.jpg')" }}
	>
		<div className="overlay"></div>
		<Container>
			<Row>
				<Col lg="8" md="10" className="mx-auto">
					<div className="site-heading">
						<h1>Distro Blog</h1>
						<span className="subheading">
							A Blog Theme by Start Bootstrap
						</span>
					</div>
				</Col>
			</Row>
		</Container>
	</header>
);

export default Header;
