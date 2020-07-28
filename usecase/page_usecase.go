package usecase

import (
	"context"
	"github.com/reecerussell/distro-blog/domain/service"
	"net/http"
	"strings"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type PageUsecase interface {
	CreatePage(ctx context.Context, d *dto.CreatePage) result.Result
	CreateBlog(ctx context.Context, d *dto.CreatePage) result.Result
	ListPages(ctx context.Context) result.Result
	ListBlogs(ctx context.Context) result.Result
	Get(ctx context.Context, id string, expand ...string) result.Result
	Update(ctx context.Context, d *dto.UpdatePage, imageData []byte) result.Result
	Activate(ctx context.Context, id string) result.Result
	Deactivate(ctx context.Context, id string) result.Result
	Delete(ctx context.Context, id string) result.Result
	GetDropdownOptions(ctx context.Context) result.Result
}

type pageUsecase struct {
	repo repository.PageRepository
	svc *service.PageService
	media MediaUsecase
}

func NewPageUsecase(repo repository.PageRepository, media MediaUsecase) PageUsecase {
	return &pageUsecase{
		repo: repo,
		svc: service.NewPageService(repo),
		media: media,
	}
}

func (u *pageUsecase) CreatePage(ctx context.Context, d *dto.CreatePage) result.Result {
	p, err := model.NewPage(ctx, d)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	return u.createPage(ctx, p)
}

func (u *pageUsecase) CreateBlog(ctx context.Context, d *dto.CreatePage) result.Result {
	p, err := model.NewBlogPage(ctx, d)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	return u.createPage(ctx, p)
}

func (u *pageUsecase) createPage(ctx context.Context, p *model.Page) result.Result {
	res := u.svc.EnsureURLIsUnique(ctx, p)
	if !res.IsOk() {
		return res
	}

	success, status, _, err := u.repo.Create(ctx, p).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	return result.Ok().WithValue(p.GetID())
}

func (u *pageUsecase) ListPages(ctx context.Context) result.Result {
	return u.repo.ListPages(ctx)
}

func (u *pageUsecase) ListBlogs(ctx context.Context) result.Result {
	return u.repo.ListBlogs(ctx)
}

func (u *pageUsecase) Get(ctx context.Context, id string, expand ...string) result.Result {
	res := u.repo.Get(ctx, id)
	if !res.IsOk(){
		return res
	}

	success, status, value, err := res.Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	p := value.(*model.Page).DTO()

	for _, e := range expand {
		switch strings.ToLower(e) {
		case "audit":
			logging.Debugf("Expanded Audit.\n")
			success, _, audit, err := u.repo.GetAudit(ctx, id).Deconstruct()
			if success {
				p.Audit = audit.([]*dto.PageAudit)
			} else {
				logging.Errorf("An error occurred while getting the page's audit data: %v", err)
			}
		}
	}

	return result.Ok().WithValue(p)
}

func (u *pageUsecase) Update(ctx context.Context, d *dto.UpdatePage, imageData []byte) result.Result {
	logging.Debugf("Attempting to update page...\n")
	logging.Debugf("Fetching page...\n")
	success, status, value, err := u.repo.Get(ctx, d.ID).Deconstruct()
	if !success {
		logging.Errorf("Failed to fetch user: %v\n", err)
		return result.Failure(err).WithStatusCode(status)
	}

	logging.Debugf("Updating page model...\n")
	p := value.(*model.Page)
	err = p.Update(ctx, d)
	if err != nil {
		logging.Errorf("Failed to update page model: %v\n",err)
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	res := u.svc.EnsureURLIsUnique(ctx, p)
	if !res.IsOk() {
		return res
	}

	if imageData != nil {
		logging.Debugf("Updating page image...\n")
		logging.Debugf("Image size: %d\n", len(imageData))
		success, status, value, err := u.media.Upload(ctx, imageData).Deconstruct()
		if !success {
			logging.Errorf("Failed to update page image: %v\n", err)
			return result.Failure(err).WithStatusCode(status)
		}

		p.UpdateImage(ctx, value.(*model.Image))
	}

	logging.Debugf("Saving changes...\n")
	success, status, _, err = u.repo.Update(ctx, p).Deconstruct()
	if !success {
		logging.Debugf("Failed to save changes: %v\n", err)
		return result.Failure(err).WithStatusCode(status)
	}

	return result.Ok()
}

func (u *pageUsecase) Activate(ctx context.Context, id string) result.Result {
	res := u.repo.Get(ctx, id)
	if !res.IsOk(){
		return res
	}

	success, status, value, err := res.Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	p := value.(*model.Page)
	err = p.Activate(ctx)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	res = u.repo.Update(ctx, p)
	if !res.IsOk() {
		return res
	}

	return result.Ok()
}

func (u *pageUsecase) Deactivate(ctx context.Context, id string) result.Result {
	res := u.repo.Get(ctx, id)
	if !res.IsOk(){
		return res
	}

	success, status, value, err := res.Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	p := value.(*model.Page)
	err = p.Deactivate(ctx)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	res = u.repo.Update(ctx, p)
	if !res.IsOk() {
		return res
	}

	return result.Ok()
}

func (u *pageUsecase) Delete(ctx context.Context, id string) result.Result {
	logging.Debugf("Attempting to delete page...\n")
	logging.Debugf("Getting page...\n")
	res := u.repo.Get(ctx, id)
	success, status, value, err := res.Deconstruct()
	if !success {
		logging.Errorf("Failed to get page: %v\n", err)
		return result.Failure(err).WithStatusCode(status)
	}

	p := value.(*model.Page)
	if imgID := p.GetImageID(); imgID != nil {
		logging.Debugf("Deleting page image...\n")
		success, status, _, err = u.media.Delete(ctx, *imgID).Deconstruct()
		if !success {
			logging.Errorf("Failed to delete image: %v\n", err)
			return result.Failure(err).WithStatusCode(status)
		}
	}

	logging.Debugf("Deleting page...\n")
	success, status, _, err = u.repo.Delete(ctx, id).Deconstruct()
	if !success {
		logging.Errorf("Failed to delete page: %v\n", err)
		return result.Failure(err).WithStatusCode(status)
	}
	
	return result.Ok()
}

func (u *pageUsecase) GetDropdownOptions(ctx context.Context) result.Result {
	return u.repo.GetDropdownOptions(ctx)
}