package datamodel

import (
	"database/sql"
	"time"
)

type UserAudit struct {
	Message string
	Date time.Time
	UserID string
	UserFullname string
	State sql.NullString
}
