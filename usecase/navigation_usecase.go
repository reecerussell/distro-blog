package usecase

import (
	"context"
	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/domain/service"
	"github.com/reecerussell/distro-blog/libraries/result"
	"net/http"
)

type NavigationUsecase interface {
	GetItems(ctx context.Context) result.Result
	GetItem(ctx context.Context, id string) result.Result
	CreateItem(ctx context.Context, d *dto.NavigationItem) result.Result
	UpdateItem(ctx context.Context, d *dto.NavigationItem) result.Result
	DeleteItem(ctx context.Context, id string) result.Result
}

type navigationUsecase struct {
	repo repository.NavigationRepository
	svc *service.NavigationService
}

func NewNavigationUsecase(repo repository.NavigationRepository) NavigationUsecase{
	return &navigationUsecase{
		repo: repo,
		svc: service.NewNavigationService(repo),
	}
}

func (u *navigationUsecase) GetItems(ctx context.Context) result.Result {
	success, _, value, err := u.repo.GetItems(ctx).Deconstruct()
	if !success {
		return result.Failure(err)
	}

	items := value.([]*model.NavigationItem)
	dtos := make([]*dto.NavigationItem, len(items))

	for i, item := range items {
		dtos[i] = item.DTO()
	}

	return result.Ok().WithValue(dtos)
}

func (u *navigationUsecase) GetItem(ctx context.Context, id string) result.Result {
	success, status, item, err := u.repo.GetItem(ctx, id).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	dto := item.(*model.NavigationItem).DTO()
	return result.Ok().WithValue(dto)
}

func (u *navigationUsecase) CreateItem(ctx context.Context, d *dto.NavigationItem) result.Result {
	ni, err := model.NewNavigationItem(d)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	res := u.repo.CreateItem(ctx, ni)
	if !res.IsOk() {
		return res
	}

	res = u.svc.EnsureSingleBrand(ctx, ni)
	if !res.IsOk() {
		return res
	}

	return result.Ok()
}

func (u *navigationUsecase) UpdateItem(ctx context.Context, d *dto.NavigationItem) result.Result {
	success, status, value, err := u.repo.GetItem(ctx, d.ID).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	ni := value.(*model.NavigationItem)
	err = ni.Update(d)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	res := u.repo.UpdateItem(ctx, ni)
	if !res.IsOk() {
		return res
	}

	res = u.svc.EnsureSingleBrand(ctx, ni)
	if !res.IsOk() {
		return res
	}

	return result.Ok()
}

func (u *navigationUsecase) DeleteItem(ctx context.Context, id string) result.Result {
	return u.repo.DeleteItem(ctx, id)
}