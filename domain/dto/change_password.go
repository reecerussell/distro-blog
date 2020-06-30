package dto

// ChangePassword is a data-transfer object used to handle data
// when changing a user's password.
type ChangePassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword string `json:"newPassword"`
}
