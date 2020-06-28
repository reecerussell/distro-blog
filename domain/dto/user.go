package dto

// User is a data-transfer object for the user domain. Used to
// handle integrations with the API layer.
type User struct {
	ID string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Email string `json:"email"`
	NormalizedEmail string `json:"normalizedEmail"`

	Audit []*UserAudit `json:"audit,omitempty"`
}
