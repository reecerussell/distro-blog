package mysql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/logging"
	"github.com/reecerussell/distro-blog/libraries/result"
)

var (
	errMsgNavigationServerError = "NAVIGATION_SERVER_ERROR"
)

type navigationRepository struct {
	db *database.MySQL
}

func NewNavigationRepository(db *database.MySQL) repository.NavigationRepository {
	return &navigationRepository{
		db: db,
	}
}

func (r *navigationRepository) GetItems(ctx context.Context) result.Result {
	const query string = "SELECT * FROM `view_navigation_items`;"
	items, err := r.db.Multiple(ctx, query, navigationItemReader)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgNavigationServerError)
	}

	models := make([]*model.NavigationItem, len(items))

	for i, item := range items {
		models[i] = model.NavigationItemFromDataModel(item.(*datamodel.NavigationItem))
	}

	return result.Ok().WithValue(models)
}

func (r *navigationRepository) GetItem(ctx context.Context, id string) result.Result {
	const query string = "SELECT * FROM `view_navigation_items` WHERE `Id` = ?;"
	item, err := r.db.Read(ctx, query, navigationItemReader, id)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgNavigationServerError)
	}

	if item == nil {
		msg := fmt.Sprintf("No navigation item was found with id '%s'", id)
		return result.Failure(msg).WithStatusCode(http.StatusNotFound)
	}

	ni := model.NavigationItemFromDataModel(item.(*datamodel.NavigationItem))
	return result.Ok().WithValue(ni)
}

func navigationItemReader(s database.ScannerFunc) (interface{}, error) {
	var dm datamodel.NavigationItem
	err := s(
		&dm.ID,
		&dm.Text,
		&dm.Target,
		&dm.URL,
		&dm.PageID,
		&dm.IsHidden,
		&dm.IsBrand,
	)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}

func (r *navigationRepository) CreateItem(ctx context.Context, ni *model.NavigationItem) result.Result {
	const query string = "CALL `create_navigation_item`(?,?,?,?,?,?,?);"
	dm := ni.DataModel()
	args := []interface{}{
		dm.ID,
		dm.Text,
		dm.Target,
		dm.URL,
		dm.PageID,
		dm.IsHidden,
		dm.IsBrand,
	}

	_, err := r.db.Execute(ctx, query, args...)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgNavigationServerError)
	}

	return result.Ok()
}

func (r *navigationRepository) UpdateItem(ctx context.Context, ni *model.NavigationItem) result.Result{
	const query string = "CALL `update_navigation_item`(?,?,?,?,?,?,?);"
	dm := ni.DataModel()
	args := []interface{}{
		dm.ID,
		dm.Text,
		dm.Target,
		dm.URL,
		dm.PageID,
		dm.IsHidden,
		dm.IsBrand,
	}

	_, err := r.db.Execute(ctx, query, args...)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgNavigationServerError)
	}

	return result.Ok()
}

func (r *navigationRepository) DeleteItem(ctx context.Context, id string) result.Result {
	const query string = "DELETE FROM `navigation` WHERE `id` = ?;"
	_, err := r.db.Execute(ctx, query, id)
	if err != nil {
		logging.Error(err)
		return result.Failure(errMsgNavigationServerError)
	}

	return result.Ok()
}
