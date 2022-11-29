package crudsql

import (
	"fmt"
	"github.com/dekamik/dvbcrud-go/internal"
	"github.com/jmoiron/sqlx"
	"reflect"
)

// SQLRepository handles CRUD queries to a table in an SQL database.
// T is the struct type that will be mapped against the table rows.
type SQLRepository[T any] struct {
	db          *sqlx.DB
	dialect     SQLDialect
	tableName   string
	idFieldName string
}

// Create inserts the values in model into a new row in the table.
func (r SQLRepository[T]) Create(model T) error {
	fields, values, err := internal.ParseProperties(model, r.idFieldName)
	if err != nil {
		return err
	}

	sql, err := getInsertIntoStmt(r.dialect, r.tableName, fields...)
	if err != nil {
		return err
	}

	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	exec, err := stmt.Exec(values...)
	if err != nil {
		return err
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("%d rows affected by INSERT INTO statement", affected)
	}

	return nil
}

// Read fetches a row from the table whose ID matches id.
func (r SQLRepository[T]) Read(id any) (*T, error) {
	var result T
	fields, err := internal.ParseFieldNames(reflect.TypeOf(result))
	if err != nil {
		return nil, err
	}

	sql, err := getSelectFromStmt(r.dialect, r.tableName, r.idFieldName, fields...)
	if err != nil {
		return nil, err
	}

	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRowx(id).StructScan(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ReadAll fetches all rows from the table.
func (r SQLRepository[T]) ReadAll() ([]T, error) {
	var result []T
	fields, err := internal.ParseFieldNames(reflect.TypeOf(result).Elem())
	if err != nil {
		return nil, err
	}

	sql := getSelectAllStmt(r.tableName, fields...)
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Select(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates the row in the table, whose ID matches id, with the data found in model.
func (r SQLRepository[T]) Update(id any, model T) error {
	fields, values, err := internal.ParseProperties(model, r.idFieldName)
	if err != nil {
		return err
	}

	sql, err := getUpdateStmt(r.dialect, r.tableName, r.idFieldName, fields...)
	if err != nil {
		return err
	}

	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	allValues := append(values, id)
	exec, err := stmt.Exec(allValues...)
	if err != nil {
		return err
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("%d rows affected by UPDATE statement", affected)
	}

	return nil
}

// Delete removes the row whose ID matches id.
func (r SQLRepository[T]) Delete(id any) error {
	sql, err := getDeleteFromStmt(r.dialect, r.tableName, r.idFieldName)
	if err != nil {
		return err
	}

	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	exec, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if affected > 1 {
		return fmt.Errorf("%d rows affected by DELETE statement", affected)
	}

	return nil
}

// NewSQLRepository creates and returns a new SQLRepository.
func NewSQLRepository[T any](db *sqlx.DB, dialect SQLDialect, tableName string, idFieldName string) (*SQLRepository[T], error) {
	if db == nil {
		return nil, fmt.Errorf("db cannot be nil")
	}
	if tableName == "" {
		return nil, fmt.Errorf("tableName cannot be empty")
	}
	if idFieldName == "" {
		idFieldName = "id"
	}

	return &SQLRepository[T]{
		db:          db,
		dialect:     dialect,
		tableName:   tableName,
		idFieldName: idFieldName,
	}, nil
}
