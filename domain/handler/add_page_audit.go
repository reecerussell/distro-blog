package handler

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/event"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

// AddPageAudit is a basic type used to implement a domain event
// handler. This is used to add page audit records.
type AddPageAudit struct {}

// Invoke invokes the handler for event.AddPageAudit domain events and
// executes a stored procedure to add a page audit record.
func (*AddPageAudit) Invoke(ctx context.Context, tx *database.Transaction, e interface{}) result.Result {
	const query string = "CALL `add_page_audit`(?,?,?,?);"
	evt := e.(*event.AddPageAudit)
	args := []interface{}{
		evt.PageID,
		evt.UserID,
		evt.Date,
		evt.Message,
	}

	err := tx.Execute(ctx, query, args...)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok()
}
