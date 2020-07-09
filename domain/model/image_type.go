package model

import "github.com/reecerussell/distro-blog/domain/datamodel"

type ImageType struct {
	id string
	name string
}

func (t *ImageType) GetID() string {
	return t.id
}

func ImageTypeFromDataModel(d *datamodel.ImageType) *ImageType {
	return &ImageType{
		id: d.ID,
		name: d.Name,
	}
}