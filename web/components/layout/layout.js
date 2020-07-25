import { Navigation, Header, Footer, Seo } from "./index";
import { Container, Row, Col } from "reactstrap";
import PropTypes from "prop-types";

const Layout = ({
	children,
	isBlog,
	toggleNavigation,
	isNavigationOpen,
	data,
}) => {
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
			<Seo imageId={data.imageId} isBlog={data.isBlog} {...data.seo} />
			<Navigation toggle={toggleNavigation} isOpen={isNavigationOpen} />
			<Header
				title={data.title}
				description={data.description}
				imageId={data.imageId}
			/>
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
