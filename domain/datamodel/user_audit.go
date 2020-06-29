package datamodel

import "time"

type UserAudit struct {
	Message string
	Date time.Time
	UserID string
	UserFullname string
	State string
}
