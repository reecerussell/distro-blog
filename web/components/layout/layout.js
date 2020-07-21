import { Navigation, Header, Footer } from "./index";
import { Container, Row, Col } from "reactstrap";
import PropTypes from "prop-types";

const Layout = ({ children, isBlog, toggleNavigation, isNavigationOpen }) => {
	const content = (
		<Container>
			<Row>
				<Col lg="8" md="10" className="mx-auto">
					{children}
				</Col>
			</Row>
		</Container>
	);

	return (
		<>
			<Navigation toggle={toggleNavigation} isOpen={isNavigationOpen} />
			<Header />
			{isBlog ? <article>{content}</article> : content}
			<hr />
			<Footer />
		</>
	);
};

Layout.propTypes = {
	children: PropTypes.node,
	isBlog: PropTypes.bool,
	toggleNavigation: PropTypes.func.isRequired,
	isNavigationOpen: PropTypes.bool.isRequired,
};
Layout.defaultProps = {
	children: null,
	isBlog: false,
};

export default Layout;
