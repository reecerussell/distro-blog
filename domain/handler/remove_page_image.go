package handler

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/event"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type RemovePageImage struct {}

func (*RemovePageImage) Invoke(ctx context.Context, tx *database.Transaction, e interface{}) result.Result {
	const query string = "CALL `remove_page_image`(?);"
	evt := e.(*event.RemovePageImage)

	err := tx.Execute(ctx, query, evt.PageID)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok()
}