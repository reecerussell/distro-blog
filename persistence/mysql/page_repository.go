package mysql

import (
	"context"
	"database/sql"
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
	if err != nil && err != sql.ErrNoRows {
		logging.Error(err)
		return result.Failure(errMsgPageDbError)
	}

	if err == sql.ErrNoRows || dm == nil {
		return result.Failure(errMsgPageNotFound).WithStatusCode(http.StatusNotFound)
	}

	p := model.PageFromDataModel(dm.(*datamodel.Page))
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
	const query string = "CALL `create_page`(?,?,?,?,?);"
	dm := p.DataModel()
	args := []interface{}{
		dm.ID,
		dm.Title,
		dm.Description,
		dm.Content,
		dm.IsBlog,
	}

	return r.executePage(ctx, p, query, args)
}

func (r *pageRepository) Update(ctx context.Context, p *model.Page) result.Result {
	const query string = "CALL `update_page`(?,?,?,?,?,?);"
	dm := p.DataModel()
	args := []interface{}{
		dm.ID,
		dm.Title,
		dm.Description,
		dm.Content,
		dm.IsActive,
		dm.ImageID,
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