package datamodel

import "database/sql"

type Image struct {
	ID string
	TypeID string
	AlternativeText sql.NullString
}
