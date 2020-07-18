package datamodel

import "database/sql"

type Setting struct {
	Key string
	Value sql.NullString
}
