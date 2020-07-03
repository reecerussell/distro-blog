package event

import "time"

// AddPageAudit is a domain event object used
// to add an audit record to a page.
type AddPageAudit struct {
	PageID string
	UserID string
	Date time.Time
	Message string
}
