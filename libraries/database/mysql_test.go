package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
)

var (
	testConnString = os.Getenv("CONN_STRING")
)

/*

Many of these tests may seem/be duplicates, replicated for each function,
however, they are there to cover cases in each method and make it easier
to adapt method specific tests in future.

*/

// EXECUTE

// Covers ensureConnected, where sql.Open errors.
func TestExecuteWithInvalidConnString(t *testing.T) {
	db := NewMySQL("invalid connection")

	ctx := context.Background()
	_, err := db.Execute(ctx, "query should be needed")
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

// Covers ensureConnected, where db.Ping returns an error.
func TestExecuteWithClosedConnection(t *testing.T) {
	db := NewMySQL(testConnString)

	err := db.ensureConnected()
	if err != nil {
		panic(err)
	}

	// Close connection
	db.db.Close()

	ctx := context.Background()
	_, err = db.Execute(ctx, "shouldn't be needed")
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestExecuteInsert(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "INSERT INTO `table-one` (`name`, `age`) VALUES (?,?);"
	args := []interface{}{"John", 20}

	ra, err := db.Execute(ctx, query, args...)
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}

	if ra != 1 {
		t.Errorf("expected to affect 1 row, but affected: %d", ra)
	}

	// clean up
	executeHelper("delete from `table-one` where `name` = ?;", args[0])
}

func TestExecuteInvalidQuery(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "invalid query"

	_, err := db.Execute(ctx, query)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestExecuteWithInvalidNumOfArgs(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "insert into `table-one` (`name`, `age`) values (?,?);"

	// too few args
	args := []interface{}{1}
	_, err := db.Execute(ctx, query, args...)
	if err == nil {
		t.Errorf("few: expected error but got none")
	}

	// too many args
	args = []interface{}{1, 2, 3}
	_, err = db.Execute(ctx, query, args...)
	if err == nil {
		t.Errorf("many: expected error but got none")
	}
}

// READ

// Covers ensureConnected, where sql.Open errors.
func TestReadWithInvalidConnString(t *testing.T) {
	db := NewMySQL("invalid connection")

	ctx := context.Background()
	_, err := db.Read(ctx, "query should be needed", testReader)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

// Covers ensureConnected, where db.Ping returns an error.
func TestReadWithClosedConnection(t *testing.T) {
	db := NewMySQL(testConnString)

	err := db.ensureConnected()
	if err != nil {
		panic(err)
	}

	// Close connection
	db.db.Close()

	ctx := context.Background()
	_, err = db.Read(ctx, "shouldn't be needed", testReader)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestRead(t *testing.T) {
	// seed db
	executeHelper("insert into `table-one` (`name`,`age`) values (?,?);", "John", 2)

	db := NewMySQL(testConnString)

	ctx := context.Background()
	res, err := db.Read(ctx, "select `name`, `age` from `table-one` where `name` = ?;", testReader, "John")
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}

	// test reader
	if res.(string) != "John, 2" {
		t.Errorf("invalid data: %v", res)
	}

	// clean up
	executeHelper("delete from `table-one` where `name` = ?;", "John")
}

func TestReadWithInvalidQuery(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	_, err := db.Read(ctx, "invalid query", testReader)
	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestReadWithInvalidNumOfArgs(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "select * from `table-one` where id = ?;"

	// too few args
	args := []interface{}{}
	_, err := db.Read(ctx, query, testReader, args...)
	if err == nil {
		t.Errorf("few: expected error but got none")
	}

	// too many args
	args = []interface{}{1, 2}
	_, err = db.Read(ctx, query, testReader, args...)
	if err == nil {
		t.Errorf("many: expected error but got none")
	}
}

func TestReadWithInvalidReader(t *testing.T) {
	// seed db
	executeHelper("insert into `table-one` (`name`,`age`) values (?,?);", "John", 2)

	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "select `name`, `age` from `table-one` where `name` = ?;"
	args := []interface{}{"John"}

	_, err := db.Read(ctx, query, testInvalidReader, args...)
	if err == nil {
		t.Errorf("expected an error, but got nil")
	}

	// clean up
	executeHelper("delete from `table-one` where `name` = ?;", "John")
}

func TestReadWithNoResults(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "select `name`, `age` from `table-one` where `name` = ?;"
	args := []interface{}{"some random name"}

	res, err := db.Read(ctx, query, testReader, args...)
	if err != nil {
		t.Errorf("expected nil, but got: %v", err)
	}

	if res != nil {
		t.Errorf("expected a nil-value, but got: %v", res)
	}
}

// MULTIPLE

// Covers ensureConnected, where sql.Open errors.
func TestMultipleWithInvalidConnString(t *testing.T) {
	db := NewMySQL("invalid connection")

	ctx := context.Background()
	_, err := db.Multiple(ctx, "query should be needed", testReader)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

// Covers ensureConnected, where db.Ping returns an error.
func TestMultipleWithClosedConnection(t *testing.T) {
	db := NewMySQL(testConnString)

	err := db.ensureConnected()
	if err != nil {
		panic(err)
	}

	// Close connection
	db.db.Close()

	ctx := context.Background()
	_, err = db.Multiple(ctx, "shouldn't be needed", testReader)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestMultiple(t *testing.T) {
	// seed db
	executeHelper("insert into `table-one` (`name`,`age`) values (?,?);", "John", 2)

	db := NewMySQL(testConnString)

	ctx := context.Background()
	results, err := db.Multiple(ctx, "select `name`, `age` from `table-one`;", testReader)
	if err != nil {
		t.Errorf("expected nil but got: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, but got: %d", len(results))
	} else {
		res := results[0]
		// test reader
		if res.(string) != "John, 2" {
			t.Errorf("invalid data: %v", res)
		}
	}

	// clean up
	executeHelper("delete from `table-one` where `name` = ?;", "John")
}

func TestMultipleWithInvalidQuery(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	_, err := db.Multiple(ctx, "invalid query", testReader)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestMultipleWithInvalidNumOfArgs(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "select * from `table-one` where id = ?;"

	// too few args
	args := []interface{}{}
	_, err := db.Multiple(ctx, query, testReader, args...)
	if err == nil {
		t.Errorf("few: expected error but got none")
	}

	// too many args
	args = []interface{}{1, 2}
	_, err = db.Multiple(ctx, query, testReader, args...)
	if err == nil {
		t.Errorf("many: expected error but got none")
	}
}

func TestMultipleWithInvalidReader(t *testing.T) {
	// seed db
	executeHelper("insert into `table-one` (`name`,`age`) values (?,?);", "John", 2)

	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "select `name`, `age` from `table-one`;"

	_, err := db.Multiple(ctx, query, testInvalidReader)
	if err == nil {
		t.Errorf("expected an error, but got nil")
	}

	// clean up
	executeHelper("delete from `table-one` where `name` = ?;", "John")
}

// COUNT

// Covers ensureConnected, where sql.Open errors.
func TestCountWithInvalidConnString(t *testing.T) {
	db := NewMySQL("invalid connection")

	ctx := context.Background()
	_, err := db.Multiple(ctx, "query should be needed", testReader)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

// Covers ensureConnected, where db.Ping returns an error.
func TestCountWithClosedConnection(t *testing.T) {
	db := NewMySQL(testConnString)

	err := db.ensureConnected()
	if err != nil {
		panic(err)
	}

	// Close connection
	db.db.Close()

	ctx := context.Background()
	_, err = db.Count(ctx, "shouldn't be needed", testReader)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestCount(t *testing.T) {
	// seed db
	executeHelper("insert into `table-one` (`name`,`age`) values (?,?);", "John", 2)
	executeHelper("insert into `table-one` (`name`,`age`) values (?,?);", "Jane", 5)

	db := NewMySQL(testConnString)
	ctx := context.Background()
	query := "SELECT COUNT(*) FROM `table-one`;"

	c, err := db.Count(ctx, query)
	if err != nil {
		t.Errorf("expected nil, but got: %v", err)
	}

	if c != 2 {
		t.Errorf("expected count to result in 2 but got: %d", c)
	}

	// clean up
	executeHelper("delete from `table-one`;")
}

func TestCountWithNonIntValue(t *testing.T) {
	// seed db
	executeHelper("insert into `table-one` (`name`,`age`) values (?,?);", "John", 2)
	executeHelper("insert into `table-one` (`name`,`age`) values (?,?);", "Jane", 5)

	db := NewMySQL(testConnString)
	ctx := context.Background()
	query := "SELECT 'hello';"

	_, err := db.Count(ctx, query)
	if err == nil {
		t.Errorf("expected error, but got nil")
	}

	// clean up
	executeHelper("delete from `table-one`;")
}

func TestCountWithInvalidQuery(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	_, err := db.Count(ctx, "invalid query")
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestCountWithInvalidNumOfArgs(t *testing.T) {
	db := NewMySQL(testConnString)

	ctx := context.Background()
	query := "select count(*) from `table-one` where id = ?;"

	// too few args
	args := []interface{}{}
	_, err := db.Count(ctx, query, args...)
	if err == nil {
		t.Errorf("few: expected error but got none")
	}

	// too many args
	args = []interface{}{1, 2}
	_, err = db.Count(ctx, query, args...)
	if err == nil {
		t.Errorf("many: expected error but got none")
	}
}

// HELPERS

func executeHelper(query string, args ...interface{}) {
	db, err := sql.Open("mysql", testConnString)
	if err != nil {
		panic(fmt.Errorf("open: %v", err))
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		panic(fmt.Errorf("exec: %v", err))
	}
}

// TEST READERS

func testReader(s ScannerFunc) (interface{}, error) {
	var (
		name string
		age  int
	)

	err := s(&name, &age)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%s, %d", name, age), err
}

func testInvalidReader(s ScannerFunc) (interface{}, error) {
	var (
		name int
		age  string
	)

	err := s(&name, &age)
	if err != nil {
		return nil, err
	}

	// should never get this far.
	return fmt.Sprintf("%d, %s", name, age), err
}
