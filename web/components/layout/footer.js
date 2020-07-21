import React from "react";
import { Container, Row, Col } from "reactstrap";

const Footer = () => (
	<footer>
		<Container>
			<Row>
				<Col lg="8" md="10" className="mx-auto">
					<ul className="list-inline text-center">
						<li className="list-inline-item">
							<a href="#">
								<span className="fa-stack fa-lg">
									<i className="fas fa-circle fa-stack-2x"></i>
									<i className="fab fa-twitter fa-stack-1x fa-inverse"></i>
								</span>
							</a>
						</li>
						<li className="list-inline-item">
							<a href="#">
								<span className="fa-stack fa-lg">
									<i className="fas fa-circle fa-stack-2x"></i>
									<i className="fab fa-twitter fa-stack-1x fa-inverse"></i>
								</span>
							</a>
						</li>
						<li className="list-inline-item">
							<a href="#">
								<span className="fa-stack fa-lg">
									<i className="fas fa-circle fa-stack-2x"></i>
									<i className="fab fa-twitter fa-stack-1x fa-inverse"></i>
								</span>
							</a>
						</li>
					</ul>
					<p className="copyright text-muted">
						Copyright &copy; Your Website 2020
					</p>
				</Col>
			</Row>
		</Container>
	</footer>
);

export default Footer;
