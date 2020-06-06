package datamodel

// User is a datamodel for the User domain.
type User struct {
	ID              string
	Firstname       string
	Lastname        string
	Email           string
	NormalizedEmail string
	PasswordHash    string
}
