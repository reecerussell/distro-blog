package model

import (
	"database/sql"
	"fmt"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
)

type SEO struct {
	title *string
	description *string
	index bool
	follow bool
}

func NewSEO(d *dto.SEO) (*SEO, error) {
	var seo SEO
	err := seo.Update(d)
	return &seo, err
}

func (seo *SEO) UpdateTitle(title *string) error {
	if title != nil && len(*title) > 255 {
		return fmt.Errorf("title cannot be greater than 255 characters long")
	}

	if title != nil && *title == "" {
		title = nil
	}

	seo.title = title

	return nil
}

func (seo *SEO) UpdateDescription(description *string) error {
	if description != nil && len(*description) > 255 {
		return fmt.Errorf("description cannot be greater than 255 characters long")
	}

	if description != nil && *description == "" {
		description = nil
	}

	seo.description = description

	return nil
}

func (seo *SEO) ShouldIndex(index bool) {
	seo.index = index
}

func (seo *SEO) ShouldFollow(follow bool) {
	seo.follow = follow
}

func (seo *SEO) Update(d *dto.SEO) error {
	err := seo.UpdateTitle(d.Title)
	if err != nil {
		return err
	}

	err = seo.UpdateDescription(d.Description)
	if err != nil {
		return err
	}

	seo.ShouldIndex(d.Index)
	seo.ShouldFollow(d.Follow)

	return nil
}

func (seo *SEO) DTO() *dto.SEO {
	return &dto.SEO{
		Title: seo.title,
		Description: seo.description,
		Index: seo.index,
		Follow: seo.follow,
	}
}

func (seo *SEO) DataModel() *datamodel.SEO {
	dm := &datamodel.SEO{
		Index: seo.index,
		Follow: seo.follow,
	}

	if seo.title == nil {
		dm.Title = sql.NullString{
			Valid: false,
		}
	} else {
		dm.Title = sql.NullString{
			Valid: true,
			String: *seo.title,
		}
	}

	if seo.description == nil {
		dm.Description = sql.NullString{
			Valid: false,
		}
	} else {
		dm.Description = sql.NullString{
			Valid: true,
			String: *seo.description,
		}
	}

	return dm
}

func SEOFromDataModel(d *datamodel.SEO) *SEO {
	seo := &SEO {
		index: d.Index,
		follow: d.Follow,
	}

	if d.Title.Valid {
		seo.title = &d.Title.String
	}

	if d.Description.Valid {
		seo.description = &d.Description.String
	}

	return seo
}
