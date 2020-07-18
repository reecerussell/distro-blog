package mysql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/reecerussell/distro-blog/domain/datamodel"
	"github.com/reecerussell/distro-blog/domain/model"
	"github.com/reecerussell/distro-blog/domain/repository"
	"github.com/reecerussell/distro-blog/libraries/database"
	"github.com/reecerussell/distro-blog/libraries/result"
)

type settingRepository struct {
	db *database.MySQL
}

func NewSettingRepository(db *database.MySQL) repository.SettingRepository {
	return &settingRepository{
		db: db,
	}
}

func (r *settingRepository) List(ctx context.Context) result.Result {
	const query string = "SELECT * FROM `view_setting_list`;"
	items, err := r.db.Multiple(ctx, query, settingReaderFunc)
	if err != nil {
		return result.Failure(err)
	}

	settings := make([]*model.Setting, len(items))

	for i, dm := range items {
		settings[i] = model.SettingFromDataModel(dm.(*datamodel.Setting))
	}

	return result.Ok().WithValue(settings)
}

func (r *settingRepository) Get(ctx context.Context, key string) result.Result {
	const query string = "CALL `get_setting`(?)"
	item, err := r.db.Read(ctx, query, settingReaderFunc, key)
	if err != nil {
		return result.Failure(err)
	}

	if item == nil {
		msg := fmt.Sprintf("No setting exists with key '%s'.", key)
		return result.Failure(msg).WithStatusCode(http.StatusNotFound)
	}

	s := model.SettingFromDataModel(item.(*datamodel.Setting))
	return result.Ok().WithValue(s)
}

func settingReaderFunc(s database.ScannerFunc) (interface{}, error) {
	var dm datamodel.Setting
	err := s(
		&dm.Key,
		&dm.Value,
	)
	return dm, err
}

func (r *settingRepository) Update(ctx context.Context, d *model.Setting) result.Result {
	const query string = "CALL `update_setting`(?,?);"
	dm := d.DataModel()
	args := []interface{}{
		dm.Key,
		dm.Value,
	}

	_, err := r.db.Execute(ctx, query, args...)
	if err != nil {
		return result.Failure(err)
	}

	return result.Ok()
}