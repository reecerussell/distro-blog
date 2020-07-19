package datamodel

import "database/sql"

type SEO struct {
	Title sql.NullString
	Description sql.NullString
	Index bool
	Follow bool
}
