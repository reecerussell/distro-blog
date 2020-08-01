package datamodel

import "database/sql"

type NavigationItem struct {
	ID string
	Text string
	Target string
	URL sql.NullString
	PageID sql.NullString
	IsHidden bool
	IsBrand bool
}
