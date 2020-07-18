package model

import (
	"database/sql"
	"fmt"
	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
)

const (
	SettingSiteName = "SITE_NAME"
	SettingTitleFormat = "TITLE_FORMAT"
)

type Setting struct {
	key string
	value *string
}

func (s *Setting) Update(value *string) error {
	switch s.key {
	case SettingSiteName:
		return s.updateSiteName(value)
	case SettingTitleFormat:
		return s.updateTitleFormat(value)
	default:
		s.value = value
		return nil
	}
}

func (s *Setting) updateSiteName(name *string) error {
	if name == nil || *name == "" {
		return fmt.Errorf("site name can not be empty")
	}

	if len(*name) > 255 {
		return fmt.Errorf("site name cannot be greater than 255 characters long")
	}

	s.value = name

	return nil
}

func (s *Setting) updateTitleFormat(format *string) error {
	if format == nil || *format == "" {
		return fmt.Errorf("title format can not be empty")
	}

	if len(*format) > 255 {
		return fmt.Errorf("title format cannot be greater than 255 characters long")
	}

	s.value = format

	return nil
}

func (s *Setting) DTO() *dto.Setting {
	return &dto.Setting{
		Key: s.key,
		Value: s.value,
	}
}

func (s *Setting) DataModel() *datamodel.Setting {
	dm := &datamodel.Setting{
		Key: s.key,
	}

	if s.value == nil {
		dm.Value = sql.NullString{
			Valid:  false,
		}
	} else {
		dm.Value = sql.NullString{
			Valid: true,
			String: *s.value,
		}
	}

	return dm
}

func SettingFromDataModel(d *datamodel.Setting) *Setting {
	s := &Setting{
		key: d.Key,
	}

	if d.Value.Valid {
		s.value = &d.Value.String
	}

	return s
}