package dto

// PageListItem is a data transfer object used for retrieving
// page records for a list, from a data source.
type PageListItem struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
}
