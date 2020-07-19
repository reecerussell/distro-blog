package datamodel

import "database/sql"

// Page is a data model used to read and write data
// from a data source.
type Page struct {
	ID string
	Title string
	Description string
	ImageID sql.NullString
	Content sql.NullString
	IsBlog bool
	IsActive bool
	URL string

	Seo *SEO
}