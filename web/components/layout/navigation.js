import React from "react";
import {
	Collapse,
	Navbar,
	NavbarToggler,
	NavbarBrand,
	Nav,
	NavItem,
	NavLink,
	Container,
} from "reactstrap";
import PropTypes from "prop-types";

const Navigation = ({ isOpen, toggle }) => (
	<Navbar expand="lg" light fixed="top" id="mainNav">
		<Container>
			<NavbarBrand href="/">Distro Blog</NavbarBrand>
			<NavbarToggler onClick={toggle} />
			<Collapse isOpen={isOpen} navbar>
				<Nav className="ml-auto" navbar>
					<NavItem>
						<NavLink href="/">Home</NavLink>
					</NavItem>
					<NavItem>
						<NavLink href="/">About</NavLink>
					</NavItem>
					<NavItem>
						<NavLink href="/">Sample Post</NavLink>
					</NavItem>
					<NavItem>
						<NavLink href="/">Contact</NavLink>
					</NavItem>
				</Nav>
			</Collapse>
		</Container>
	</Navbar>
);

Navigation.propTypes = {
	isOpen: PropTypes.bool.isRequired,
	toggle: PropTypes.func.isRequired,
};

export default Navigation;
