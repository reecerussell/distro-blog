package model

import (
	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
)

type SEO struct {
	title *string
	description *string
	index bool
	follow bool
}

func (seo *SEO) DTO() *dto.SEO {
	return &dto.SEO{
		Title: seo.title,
		Description: seo.description,
		Index: seo.index,
		Follow: seo.follow,
	}
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
