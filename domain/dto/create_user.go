package dto

// CreateUser is a data transfer object used receive data
// needed to create a new user.
type CreateUser struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
