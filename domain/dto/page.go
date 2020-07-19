package dto

// Page is a data transfer object for the page domain.
type Page struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Content *string `json:"content"`
	IsBlog bool `json:"isBlog"`
	IsActive bool `json:"isActive"`
	ImageID *string `json:"imageId"`
	URL string `json:"url"`

	Audit []*PageAudit `json:"audit,omitempty"`
	SEO *SEO `json:"seo,omitempty"`
}
