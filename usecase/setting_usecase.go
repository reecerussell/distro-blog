package usecase

import (
	"context"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type SettingUsecase interface {
	List(ctx context.Context) result.Result
	Get(ctx context.Context, key string) result.Result
	Update(ctx context.Context, d *dto.Setting) result.Result
}

type settingUsecase struct {
	repo repository.SettingRepository
}

func NewSettingUsecase(repo repository.SettingRepository) SettingUsecase {
	return &settingUsecase{
		repo: repo,
	}
}

func (u *settingUsecase) List(ctx context.Context) result.Result {
	success, _, value, err := u.repo.List(ctx).Deconstruct()
	if !success {
		return result.Failure(err)
	}

	settings := value.([]*model.Setting)
	dtos := make([]*dto.Setting, len(settings))

	for i, s := range settings {
		dtos[i] = s.DTO()
	}

	return result.Ok().WithValue(dtos)
}

func (u *settingUsecase) Get(ctx context.Context, key string) result.Result {
	success, status, value, err := u.repo.Get(ctx, d.Key).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	s := value.(*model.Setting)
	return result.Ok().WithValue(s.DTO())
}

func (u *settingUsecase) Update(ctx context.Context, d *dto.Setting) result.Result {
	success, status, value, err := u.repo.Get(ctx, d.Key).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	s := value.(*model.Setting)
	err = s.Update(d.Value)
	if err != nil {
		return result.Failure(err).WithStatusCode(http.StatusBadRequest)
	}

	res := u.repo.Update(ctx, s)
	if !res.IsOk() {
		return res
	}

	return result.Ok()
}
