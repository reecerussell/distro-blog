package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
)

type NavigationItem struct {
	id string
	text string
	target string
	url *string
	pageID *string
	isHidden bool
	isBrand bool
}

func NewNavigationItem(d *dto.NavigationItem) (*NavigationItem, error) {
	ni := &NavigationItem{
		id: uuid.New().String(),
	}

	err := ni.Update(d)
	if err != nil {
		return nil, err
	}

	return ni, nil
}

func (ni *NavigationItem) ID() string {
	return ni.id
}

func (ni *NavigationItem) IsBrand() bool {
	return ni.isBrand
}

func (ni *NavigationItem) Update(d *dto.NavigationItem) error {
	err := ni.UpdateText(d.Text)
	if err != nil {
		return err
	}

	err = ni.UpdateURL(d.URL, d.PageID)
	if err != nil {
		return err
	}

	err = ni.UpdateTarget(d.Target)
	if err != nil {
		return err
	}

	ni.UpdateVisibility(d.IsHidden)
	ni.UpdateIsBrand(d.IsBrand)

	return nil
}

func (ni *NavigationItem) UpdateText(text string) error {
	l := len(text)
	if l == 0 {
		return fmt.Errorf("text cannot be empty")
	}

	if l > 255 {
		return fmt.Errorf("text cannot be greater than 255 characters long")
	}

	ni.text = text

	return nil
}

var allowedNavTargets = [...]string{"_self", "_blank", "_top", "_parent"}

func (ni *NavigationItem) UpdateTarget(target string) error {
	lt := strings.ToLower(target)
	found := false

	for _, s := range allowedNavTargets {
		if lt == s {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("target '%s' is not valid", target)
	}

	ni.target = lt

	return nil
}

func (ni *NavigationItem) UpdateURL(url *string, pageID *string) error {
	if (url == nil || *url == "" ) && (pageID == nil || *pageID == "") {
		return fmt.Errorf("url cannot be empty")
	}

	if pageID == nil {
		if len(*url) > 255 {
			return fmt.Errorf("url cannot be greater than 255 characters long")
		}

		ni.url = url
		ni.pageID = nil
	} else {
		if len(*pageID) > 255 {
			return fmt.Errorf("page ID cannot be greater than 128 characters long")
		}

		ni.url = nil
		ni.pageID = pageID
	}

	return nil
}

func (ni *NavigationItem) UpdateVisibility(isHidden bool) {
	ni.isHidden = isHidden
}

func (ni *NavigationItem) UpdateIsBrand(isBrand bool) {
	ni.isBrand = isBrand
}

func (ni *NavigationItem) DTO() *dto.NavigationItem {
	return &dto.NavigationItem{
		ID: ni.id,
		Text: ni.text,
		Target: ni.target,
		URL: ni.url,
		PageID: ni.pageID,
		IsBrand: ni.isBrand,
		IsHidden: ni.isHidden,
	}
}

func (ni *NavigationItem) DataModel() *datamodel.NavigationItem {
	dm := &datamodel.NavigationItem{
		ID: ni.id,
		Text: ni.text,
		Target: ni.target,
		IsBrand: ni.isBrand,
		IsHidden: ni.isHidden,
	}

	if ni.url == nil {
		dm.URL = sql.NullString{
			Valid: false,
		}
	} else {
		dm.URL = sql.NullString{
			Valid: true,
			String: *ni.url,
		}
	}

	if ni.pageID == nil {
		dm.PageID = sql.NullString{
			Valid: false,
		}
	} else {
		dm.PageID = sql.NullString{
			Valid: true,
			String: *ni.pageID,
		}
	}

	return dm
}

func NavigationItemFromDataModel(dm *datamodel.NavigationItem) *NavigationItem {
	ni := &NavigationItem{
		id: dm.ID,
		text: dm.Text,
		target: dm.Target,
		isHidden: dm.IsHidden,
		isBrand: dm.IsBrand,
	}

	if dm.URL.Valid {
		ni.url = &dm.URL.String
	}

	if dm.PageID.Valid {
		ni.pageID = &dm.PageID.String
	}

	return ni
}