package dto

// CreatePage is a data-transfer object used to read data
// from request bodies, then create page records.
type CreatePage struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Content *string `json:"content"`
	URL string `json:"url"`
}
