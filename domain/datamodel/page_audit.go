package datamodel

import (
	"time"
)

type PageAudit struct {
	Message string
	Date time.Time
	UserID string
	UserFullname string
}

