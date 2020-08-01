package service

import (
	"context"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type NavigationService struct {
	repo repository.NavigationRepository
}

func NewNavigationService(repo repository.NavigationRepository) *NavigationService{
	return &NavigationService{
		repo: repo,
	}
}

// EnsureSingleBrand updates any navigation items in the database if ni.isBrand
// is true. This ensures there is only one item flagged as a brand item.
func (s *NavigationService) EnsureSingleBrand(ctx context.Context, ni *model.NavigationItem) result.Result {
	if !ni.IsBrand() {
		return result.Ok()
	}

	return s.repo.UpdateBrand(ctx, ni.ID())
}