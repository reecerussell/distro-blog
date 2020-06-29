package handler

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/event"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type AddUserAudit struct {}

func (*AddUserAudit) Invoke(ctx context.Context, tx *database.Transaction, e interface{}) result.Result {
	const query string = "CALL `add_user_audit`(?,?,?,?);"
	evt := e.(*event.AddUserAudit)
	args := []interface{}{
		evt.Message,
		evt.Date,
		evt.UserID,
		evt.PerformingUserID,
	}

	err := tx.Execute(ctx, query, args...)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok()
}
