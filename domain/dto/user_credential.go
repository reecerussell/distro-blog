package dto

// UserCredential is used to retrieve user credentials and validate a user.
type UserCredential struct {
	Email string `json:"email"`
	Password string `json:"password""`
}