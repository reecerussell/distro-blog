package dto

type SEO struct {
	Title *string `json:"title"`
	Description *string `json:"description"`
	Index bool `json:"index"`
	Follow bool `json:"follow"`
}
