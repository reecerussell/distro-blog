package model

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/reecerussell/distro-blog/domain/datamodel"
)

type Image struct {
	id string
	typeID string
	alternativeText *string
}

func NewImage(t *ImageType) *Image {
	return &Image{
		id: uuid.New().String(),
		typeID: t.GetID(),
	}
}

func (i *Image) GetID() string {
	return i.id
}

func (i *Image) DataModel() *datamodel.Image {
	dm := &datamodel.Image{
		ID:              i.id,
		TypeID:          i.typeID,
	}

	if i.alternativeText == nil {
		dm.AlternativeText = sql.NullString{
			Valid: false,
		}
	} else {
		dm.AlternativeText = sql.NullString{
			Valid: true,
			String: *i.alternativeText,
		}
	}

	return dm
}

func ImageFromDataModel(d *datamodel.Image) *Image {
	i := &Image {
		id: d.ID,
		typeID: d.TypeID,
	}

	if d.AlternativeText.Valid {
		i.alternativeText = &d.AlternativeText.String
	}

	return i
}