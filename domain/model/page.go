package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/event"
	"github.com/reecerussell/distro-blog/domain/handler"
	"github.com/reecerussell/distro-blog/libraries/contextkey"
	"github.com/reecerussell/distro-blog/libraries/domainevents"
	"github.com/reecerussell/distro-blog/libraries/logging"
)

func init() {
	domainevents.RegisterEventHandler(&event.AddPageAudit{}, &handler.AddPageAudit{})
	domainevents.RegisterEventHandler(&event.RemovePageImage{}, handler.NewRemovePageImageHandler())
	domainevents.RegisterEventHandler(&event.UpdatePageSEO{}, &handler.UpdatePageSEO{})
}

const (
	AuditPageCreated = "PAGE_CREATED"
	AuditPageUpdated = "PAGE_UPDATED"
	AuditPageDeactivated = "PAGE_DEACTIVATED"
	AuditPageActivated = "PAGE_ACTIVATED"
	AuditPageImageUpdated = "PAGE_IMAGE_UPDATED"
)

// These tags are tags which are not allowed to be used
// in the page content. Page content can contain HTML, which
// therefore is not allowed to contain these HTML tags.
var disallowedContentTags = [...]string{"script","link","body","html","head"}

type Page struct {
	domainevents.Aggregate

	id string
	title string
	description string
	imageID *string
	content *string
	isBlog bool
	isActive bool
	url string
	seoID *string

	seo *SEO
}

// NewPage creates a new page domain object with the given date.
// This page will be marked as deactivated and as a not a blog.
func NewPage(ctx context.Context, d *dto.CreatePage) (*Page, error) {
	p := &Page {
		id: uuid.New().String(),
		isActive: false,
		isBlog: false,
	}

	err := p.updateContent(d.Title, d.Description, d.Content, d.URL, d.SEO)
	if err != nil {
		return nil, err
	}

	p.addAudit(ctx, AuditPageCreated)

	return p, nil
}

// NewBlogPage creates a new page which is marked as a blog. The new page
// will also be marked as deactivated.
func NewBlogPage(ctx context.Context, d *dto.CreatePage) (*Page, error) {
	p := &Page {
		id: uuid.New().String(),
		isActive: false,
		isBlog: true,
	}

	err := p.updateContent(d.Title, d.Description, d.Content, d.URL, d.SEO)
	if err != nil {
		return nil, err
	}

	p.addAudit(ctx, AuditPageCreated)

	return p, nil
}

// GetID returns the page's id.
func (p *Page) GetID() string {
	return p.id
}

// GetImageId return's the page's image id.
func (p *Page) GetImageID() *string {
	return p.imageID
}

// URL returns the page's url value.
func (p *Page) URL() string {
	return p.url
}

// Update updates the page's data, including; title, description and content.
func (p *Page) Update(ctx context.Context, d *dto.UpdatePage) error {
	err := p.updateContent(d.Title, d.Description, d.Content, d.URL, d.SEO)
	if err != nil {
		return err
	}

	p.addAudit(ctx, AuditPageUpdated)

	return nil
}

// updateContent moves the core update logic to a separate functions to avoid
// code duplication.
func (p *Page) updateContent(title, description string, content *string, url string, seo *dto.SEO) error {
	err := p.UpdateTitle(title)
	if err != nil {
		return err
	}

	err = p.UpdateDescription(description)
	if err != nil {
		return err
	}

	err = p.UpdateContent(content)
	if err != nil {
		return err
	}

	err = p.UpdateURL(url)
	if err != nil {
		return err
	}

	if seo != nil {
		var sdm *datamodel.SEO
		if p.seo == nil {
			s, err := NewSEO(seo)
			if err != nil {
				return err
			}

			sdm = s.DataModel()
		} else {
			err = p.seo.Update(seo)
			if err != nil {
				return err
			}

			sdm = p.seo.DataModel()
		}

		p.RaiseEvent(&event.UpdatePageSEO{
			PageID: p.id,
			SEO: sdm,
		})
	}

	return nil
}

// UpdateTitle updates the page's title, ensuring it's valid.
//
// Titles cannot be empty and cannot be greater than 255 characters long.
func (p *Page) UpdateTitle(title string) error {
	l := len(title)

	switch true {
	case l < 1:
		return fmt.Errorf("title is required")
	case p.title == title:
		return nil
	case l > 255:
		return fmt.Errorf("title cannot be greater than 255 characters long")
	}

	p.title = title

	return nil
}

// UpdateDescription updates the page's description.
//
// The description cannot be empty or greater than 255 characters long.
func (p *Page) UpdateDescription(description string) error {
	l := len(description)

	switch true {
	case l < 1:
		return fmt.Errorf("description is required")
	case p.description == description:
		return nil
	case l > 255:
		return fmt.Errorf("description cannot be greater than 255 characters long")
	}

	p.description = description

	return nil
}

// UpdateContent updates the page's HTML content. The given content
// cannot contain any HTML tags that are in disallowedContentTags.
func (p *Page) UpdateContent(content *string) error {
	if content != nil && len(*content) < 1 {
		content = nil
	}

	if content != nil {
		nc := strings.ToLower(*content)

		for _, s := range disallowedContentTags {
			if strings.Contains(nc, fmt.Sprintf("<%s", s)) {
				return fmt.Errorf("%s tags are not allowed in page content", s)
			}
		}
	}

	p.content = content

	return nil
}

// A whitelist of characters that are allowed in page urls.
var urlSafeCharMap = map[string]bool {
	"a": true,
	"b": true,
	"c": true,
	"d": true,
	"e": true,
	"f": true,
	"g": true,
	"h": true,
	"i": true,
	"j": true,
	"k": true,
	"l": true,
	"m": true,
	"n": true,
	"o": true,
	"p": true,
	"q": true,
	"r": true,
	"s": true,
	"t": true,
	"u": true,
	"v": true,
	"w": true,
	"x": true,
	"y": true,
	"z": true,
	"1": true,
	"2": true,
	"3": true,
	"4": true,
	"5": true,
	"6": true,
	"7": true,
	"8": true,
	"9": true,
	"0": true,
	"-": true,
}

// UpdateURL updates the page's url. The url can only contain
// whitelisted characters and cannot be greater than 255 chars long.
func (p *Page) UpdateURL(url string) error {
	if len(url) > 255 {
		return fmt.Errorf("page url cannot be greater than 255 characters long")
	}

	if url != "" {
		chars := strings.Split(url, "")
		for _, c := range chars {
			_, ok := urlSafeCharMap[c]
			if !ok {
				return fmt.Errorf("the character '%s' is anot allow in a page url", c)
			}
		}
	}

	p.url = url

	return nil
}

// Deactivate marks the page as inactive. An non-nil
// error will be returned if the page is already inactive.
func (p *Page) Deactivate(ctx context.Context) error {
	if !p.isActive {
		return fmt.Errorf("page is already inactive")
	}

	p.isActive = false
	p.addAudit(ctx, AuditPageDeactivated)

	return nil
}

// Activate marks the page as active. A non-nil error will
// be returned if the page is already active.
func (p *Page) Activate(ctx context.Context) error {
	if p.isActive {
		return fmt.Errorf("page ia already active")
	}

	p.isActive = true
	p.addAudit(ctx, AuditPageActivated)

	return nil
}

func (p *Page) UpdateImage(ctx context.Context, i *Image) {
	if p.imageID != nil {
		p.RaiseEvent(&event.RemovePageImage{
			ImageID: *p.imageID,
		})
	}

	if i == nil {
		p.imageID = nil
	} else {
		id := i.GetID()
		p.imageID = &id
	}

	p.addAudit(ctx, AuditPageImageUpdated)
}

func (p *Page) UpdateSEO(d *dto.SEO) error {
	return p.seo.Update(d)
}

// addAudit raises an audit domain event for the page.
func (p *Page) addAudit(ctx context.Context, message string) {
	logging.Debugf("[PAGE:%s]: raising audit domain event.\n", p.id)

	var userID string
	if uid := ctx.Value(contextkey.ContextKey("user_id")); uid == nil {
		logging.Errorf("failed to raise audit event, due to no user id being present in the context.\n")
		return
	} else {
		userID = uid.(string)
	}

	e := &event.AddPageAudit{
		PageID: p.id,
		UserID: userID,
		Date: time.Now().UTC(),
		Message: message,
	}

	p.RaiseEvent(e)
}

// DataModel returns a data model object for the page.
func (p *Page) DataModel() *datamodel.Page {
	dm := &datamodel.Page{
		ID:          p.id,
		Title:       p.title,
		Description: p.description,
		IsBlog:      p.isBlog,
		IsActive:    p.isActive,
		URL: p.url,
	}

	if p.content == nil {
		dm.Content = sql.NullString{
			Valid: false,
		}
	} else {
		dm.Content = sql.NullString{
			Valid: true,
			String: *p.content,
		}
	}

	if p.imageID == nil {
		dm.ImageID = sql.NullString{
			Valid: false,
		}
	} else {
		dm.ImageID = sql.NullString{
			Valid: true,
			String: *p.imageID,
		}
	}

	return dm
}

// PageFromDataModel returns a new instance of Page, populated with
// the data from the data model. This should only be used by repositories.
func PageFromDataModel(d *datamodel.Page) *Page {
	p := &Page{
		id: d.ID,
		title: d.Title,
		description: d.Description,
		isActive: d.IsActive,
		isBlog: d.IsBlog,
		url: d.URL,
	}

	if d.Content.Valid {
		p.content = &d.Content.String
	}

	if d.ImageID.Valid {
		p.imageID = &d.ImageID.String
	}

	if d.Seo != nil {
		p.seo = SEOFromDataModel(d.Seo)
	}

	return p
}

// DTO returns a *dto.Page for the page.
func (p *Page) DTO() *dto.Page {
	d := &dto.Page{
		ID:          p.id,
		Title:       p.title,
		Description: p.description,
		Content:     p.content,
		IsBlog:      p.isBlog,
		IsActive:    p.isActive,
		ImageID: p.imageID,
		URL: p.url,
	}

	if p.seo != nil {
		d.SEO = p.seo.DTO()
	}

	return d
}