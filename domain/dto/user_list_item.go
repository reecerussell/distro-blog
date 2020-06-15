package dto

// UserListItem is used to store data for a user, from a user list.
type UserListItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
