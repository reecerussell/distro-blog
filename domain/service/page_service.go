package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/result"
)

// PageService is used to provide the page domain with extra functionality,
// such as validation methods which don't belong to the domain layer.
type PageService struct {
	repo repository.PageRepository
}

// NewPageService returns a new instance of PageService with the given repository.
func NewPageService(repo repository.PageRepository) *PageService {
	return &PageService{
		repo: repo,
	}
}

// EnsureURLIsUnique ensures the given page's url is unique by
// ensuring it has not be previously used in the database.
func (s *PageService) EnsureURLIsUnique(ctx context.Context, p *model.Page) result.Result {
	success, _, value, err := s.repo.CountByURL(ctx, p).Deconstruct()
	if !success {
		return result.Failure(err)
	}

	count := value.(int64)
	if count > 0 {
		msg := fmt.Sprintf("The url '%s' is already being used.", p.URL())
		return result.Failure(msg).WithStatusCode(http.StatusBadRequest)
	}

	return result.Ok()
}
