package dto

// UpdatePage is a data-transfer object used to transfer and
// hold data required to update a page record.
type UpdatePage struct {
	ID string `json:"id"`
	Title string `json:"name"`
	Description string `json:"description"`
	Content *string `json:"content"`
}
