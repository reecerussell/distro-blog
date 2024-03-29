package mysql

import (
	"context"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/dto"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"github.com/reecerussell/distro-blog/libraries/result"
)

const (
	errMsgPageNotFound = "PAGE_NOT_FOUND"
	errMsgPageDbError = "PAGE_SERVER_ERROR"
	errMsgPageAuditDbError = "PAGE_AUDIT_SERVER_ERROR"
)

type pageRepository struct {
	db *database.MySQL
}

func NewPageRepository(db *database.MySQL) repository.PageRepository {
	return &pageRepository{
		db: db,
	}
}

func (r *pageRepository) Get(ctx context.Context, id string) result.Result {
	const query string = "CALL `get_page`(?);"
	dm, err := r.db.Read(ctx, query, pageReader, id)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgPageDbError)
	}

	page := dm.(*datamodel.Page)
	page.Seo, err = r.getPageSeo(ctx, id)
	if err != nil {
		logging.Errorf("Error getting page SEO: %v", err)
	}

	p := model.PageFromDataModel(page)
	return result.Ok().WithValue(p)
}

func pageReader(s database.ScannerFunc) (interface{}, error) {
	var dm datamodel.Page
	err := s(
		&dm.ID,
		&dm.Title,
		&dm.Description,
		&dm.Content,
		&dm.IsBlog,
		&dm.IsActive,
		&dm.ImageID,
		&dm.URL,
	)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}

func (r *pageRepository) getPageSeo(ctx context.Context, pageId string) (*datamodel.SEO, error) {
	const query string = "CALL `get_page_seo`(?);"
	dm, err := r.db.Read(ctx, query, seoReader, pageId)
	if err != nil {
		return nil, err
	}

	if dm == nil {
		return nil, nil
	}

	return dm.(*datamodel.SEO), nil
}

func seoReader(s database.ScannerFunc) (interface{}, error) {
	var dm datamodel.SEO
	err := s(
		&dm.Title,
		&dm.Description,
		&dm.Index,
		&dm.Follow,
	)
	return &dm, err
}

func (r *pageRepository) ListPages(ctx context.Context) result.Result {
	const query string = "SELECT * FROM `view_page_list`;"
	return r.getList(ctx, query)
}

func (r *pageRepository) ListBlogs(ctx context.Context) result.Result {
	const query string = "SELECT * FROM `view_blog_list`;"
	return r.getList(ctx, query)
}

func (r *pageRepository) getList(ctx context.Context, query string) result.Result {
	items, err := r.db.Multiple(ctx, query, listItemReader)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgPageDbError)
	}

	dtos := make([]*dto.PageListItem, len(items))

	for i, dm := range items {
		dtos[i] = dm.(*dto.PageListItem)
	}

	return result.Ok().WithValue(dtos)
}

func listItemReader(s database.ScannerFunc) (interface{}, error) {
	var dm dto.PageListItem
	err := s(
		&dm.ID,
		&dm.Title,
		&dm.Description,
	)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}

func (r *pageRepository) Create(ctx context.Context, p *model.Page) result.Result {
	const query string = "CALL `create_page`(?,?,?,?,?,?);"
	dm := p.DataModel()
	args := []interface{}{
		dm.ID,
		dm.Title,
		dm.Description,
		dm.Content,
		dm.IsBlog,
		dm.URL,
	}

	return r.executePage(ctx, p, query, args)
}

func (r *pageRepository) Update(ctx context.Context, p *model.Page) result.Result {
	const query string = "CALL `update_page`(?,?,?,?,?,?,?);"
	dm := p.DataModel()
	args := []interface{}{
		dm.ID,
		dm.Title,
		dm.Description,
		dm.Content,
		dm.IsActive,
		dm.ImageID,
		dm.URL,
	}

	return r.executePage(ctx, p, query, args)
}

func (r *pageRepository) executePage(ctx context.Context, p *model.Page, query string, args []interface{}) result.Result {
	tx, err := r.db.Tx(ctx)
	defer func() {
		tx.Finish(err)
	}()
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgPageDbError)
	}

	err = tx.Execute(ctx, query, args...)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgPageDbError)
	}

	var success bool
	var status int
	success, status, _, err = p.DispatchEvents(ctx, tx).Deconstruct()
	if !success {
		return result.Failure(err).WithStatusCode(status)
	}

	return result.Ok()
}

func (r *pageRepository) Delete(ctx context.Context, id string) result.Result {
	const query string = "DELETE FROM `pages` WHERE `id` = ?;"
	ra, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return result.Failure(errMsgPageDbError)
	}

	if ra < 1 {
		return result.Failure(errMsgPageNotFound).WithStatusCode(http.StatusNotFound)
	}

	return result.Ok()
}

func (r *pageRepository) GetAudit(ctx context.Context, id string) result.Result {
	const query string = "CALL `get_page_audit`(?);"
	items, err := r.db.Multiple(ctx, query, pageAuditReader, id)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgPageAuditDbError)
	}

	dtos := make([]*dto.PageAudit, len(items))

	for i, item := range items {
		dm := item.(*datamodel.PageAudit)
		dto := &dto.PageAudit{
			Message: dm.Message,
			Date: dm.Date,
			UserFullname: dm.UserFullname,
			UserID: dm.UserID,
		}

		dtos[i] = dto
	}

	return result.Ok().WithValue(dtos)
}

func pageAuditReader(s database.ScannerFunc) (interface{}, error) {
	var dm datamodel.PageAudit
	err := s(
		&dm.UserID,
		&dm.UserFullname,
		&dm.Message,
		&dm.Date,
	)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}

func (r *pageRepository) CountByURL(ctx context.Context, p *model.Page) result.Result {
	const query string = "CALL `count_pages_by_url`(?, ?);"

	dm := p.DataModel()
	c, err := r.db.Count(ctx, query, dm.URL, dm.ID)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok().WithValue(c)
}

// GetDropdownOptions returns a list of *dto.PageDropdownItem for each item in the database.
func (r *pageRepository) GetDropdownOptions(ctx context.Context) result.Result {
	const query string = "SELECT * FROM `view_page_dropdown_options`;"
	items, err := r.db.Multiple(ctx, query, pageOptionReader)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgPageDbError)
	}

	opts := make([]*dto.PageDropdownItem, len(items))

	for i, item := range items {
		opts[i] = item.(*dto.PageDropdownItem)
	}

	return result.Ok().WithValue(opts)
}

func pageOptionReader(s database.ScannerFunc) (interface{}, error) {
	var dto dto.PageDropdownItem
	err := s(
		&dto.ID,
		&dto.Title,
		&dto.URL,
		&dto.IsBlog,
		&dto.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}