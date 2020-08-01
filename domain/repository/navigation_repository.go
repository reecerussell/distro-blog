package repository

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type NavigationRepository interface {
	GetItems(ctx context.Context) result.Result
	GetItem(ctx context.Context, id string) result.Result
	CreateItem(ctx context.Context, ni *model.NavigationItem) result.Result
	UpdateItem(ctx context.Context, ni *model.NavigationItem) result.Result
	DeleteItem(ctx context.Context, id string) result.Result
	UpdateBrand(ctx context.Context, brandItemId string) result.Result
}
