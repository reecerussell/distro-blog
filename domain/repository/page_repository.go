package repository

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/libraries/result"
)

// PageRepository is a high-level interface used to read and write
// page data to and from a data source.
type PageRepository interface {
	Get(ctx context.Context, id string) result.Result
	ListPages(ctx context.Context) result.Result
	ListBlogs(ctx context.Context) result.Result
	Create(ctx context.Context, p *model.Page) result.Result
	Update(ctx context.Context, p *model.Page) result.Result
	Delete(ctx context.Context, id string) result.Result
	GetAudit(ctx context.Context, id string) result.Result
}