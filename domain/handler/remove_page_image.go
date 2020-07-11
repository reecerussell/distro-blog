package handler

import (
	"context"
	"github.com/reecerussell/distro-blog/libraries/storage"
	"os"

	"github.com/reecerussell/distro-blog/domain/event"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type RemovePageImage struct {
	stg *storage.Service
}

func NewRemovePageImageHandler() *RemovePageImage {
	stg, err := storage.New(os.Getenv("MEDIA_BUCKET_NAME"))
	if err != nil {
		panic(err)
	}

	return &RemovePageImage{stg:stg}
}

func (h *RemovePageImage) Invoke(ctx context.Context, tx *database.Transaction, e interface{}) result.Result {
	const query string = "CALL `delete_image`(?);"
	evt := e.(*event.RemovePageImage)

	err := tx.Execute(ctx, query, evt.ImageID)
	if err != nil {
		return result.Failure(err)
	}

	err = h.stg.Delete(evt.ImageID)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok()
}