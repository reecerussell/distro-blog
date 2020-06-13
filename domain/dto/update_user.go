package dto

// UpdateUser is a data-transfer object used to receive data
// to update a user with.
type UpdateUser struct {
	ID string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Email string `json:"email"`
}
