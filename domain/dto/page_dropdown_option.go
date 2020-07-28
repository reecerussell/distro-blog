package dto

// PageDropdownItem is used to return a list of pages
// to provide a dropdown box with options.
type PageDropdownItem struct {
	ID string `json:"id"`
	Title string `json:"title"`
	IsBlog bool `json:"isBlog"`
	IsActive bool `json:"isActive"`
}
