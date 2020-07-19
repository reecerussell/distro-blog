package handler

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/event"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type UpdatePageSEO struct {}

func (*UpdatePageSEO) Invoke(ctx context.Context, tx *database.Transaction, e interface{}) result.Result {
	const query string = "CALL `update_page_seo`(?,?,?,?,?);"
	evt := e.(*event.UpdatePageSEO)
	args := []interface{}{
		evt.PageID,
		evt.SEO.Title,
		evt.SEO.Description,
		evt.SEO.Index,
		evt.SEO.Follow,
	}

	err := tx.Execute(ctx, query, args...)
	if err != nil {
		logging.Debugf("Failed to create SEO for page: %v", err)
		return result.Failure(err)
	}

	return result.Ok()
}

