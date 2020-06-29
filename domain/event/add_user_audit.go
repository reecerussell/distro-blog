package event

import "time"

type AddUserAudit struct {
	Message string
	Date time.Time
	UserID string
	PerformingUserID string
}