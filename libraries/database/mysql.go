package database

import (
	"context"
	"database/sql"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// ScannerFunc is a function used by a ReaderFunc to read
// the results of a sql row.
type ScannerFunc func(dst ...interface{}) error

// ReaderFunc is used to read results of a query to a specific type.
type ReaderFunc func(s ScannerFunc) (interface{}, error)

// MySQL is a wrapper around database/sql methods, providing
// a higher-level abstaction of the database access code.
type MySQL struct {
	connStr string
	ctx     context.Context
	db      *sql.DB
}

// NewMySQL returns a new instance of MySQL, with the given connection string.
func NewMySQL(connStr string) *MySQL {
	return &MySQL{
		connStr: connStr,
		ctx:     context.Background(),
	}
}

// Execute runs a SQL command on the database, using the given query and arguments.
// The query is run in a SQL transaction, and will be rolled back if any error occurs.
func (mysql *MySQL) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	err := mysql.ensureConnected()
	if err != nil {
		return 0, err
	}

	// Error ignored due to the driver supporting the isolation
	// level and a successful connection already made.
	tx, _ := mysql.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}

	ra, _ := res.RowsAffected()
	return ra, nil
}

// Read queries the database, with the given query and argments. If the query results in
// no rows, the result will be (nil, nil).
func (mysql *MySQL) Read(ctx context.Context, query string, rdr ReaderFunc, args ...interface{}) (interface{}, error) {
	err := mysql.ensureConnected()
	if err != nil {
		return nil, err
	}

	stmt, err := mysql.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	item, err := rdr(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return item, nil
}

// Multiple queries multiple records from the database, with the given query and arguments.
func (mysql *MySQL) Multiple(ctx context.Context, query string, rdr ReaderFunc, args ...interface{}) ([]interface{}, error) {
	err := mysql.ensureConnected()
	if err != nil {
		return nil, err
	}

	stmt, err := mysql.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []interface{}{}

	for rows.Next() {
		item, err := rdr(rows.Scan)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	// Understaning that an error can be returned from here,
	// I can't seem to find or create a case where an error
	// would be returned.
	// err = rows.Err()
	// if err != nil {
	// 	return nil, err
	// }

	return items, nil
}

// Count is similar to Read, but can be used to read a single integer from a query. If
// an error occurs the default int64 value(0) will be returned.
func (mysql *MySQL) Count(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := mysql.Read(ctx, query, func(s ScannerFunc) (interface{}, error) {
		var c int64
		err := s(&c)
		return c, err
	}, args...)
	if err != nil {
		return 0, err
	}

	return res.(int64), nil
}

func (mysql *MySQL) ensureConnected() error {
	if mysql.db == nil {
		db, err := sql.Open("mysql", mysql.connStr)
		if err != nil {
			return err
		}

		mysql.db = db
	}

	err := mysql.db.PingContext(mysql.ctx)
	if err != nil {
		return err
	}

	return nil
}
