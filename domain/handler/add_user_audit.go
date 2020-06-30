package handler

import (
	"context"
	"encoding/json"

	"github.com/reecerussell/distro-blog/domain/event"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type AddUserAudit struct {}

func (*AddUserAudit) Invoke(ctx context.Context, tx *database.Transaction, e interface{}) result.Result {
	const query string = "CALL `add_user_audit`(?,?,?,?,?);"
	evt := e.(*event.AddUserAudit)
	state := map[string]interface{}{
		"before": evt.Before,
		"after": evt.After,
	}
	stateJson, _ := json.Marshal(state)
	args := []interface{}{
		evt.Message,
		evt.Date,
		evt.UserID,
		evt.PerformingUserID,
		string(stateJson),
	}

	err := tx.Execute(ctx, query, args...)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok()
}
