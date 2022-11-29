package dvbcrud

import (
	"fmt"
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

func (r SQLRepository[T]) Create(model T) error {
	fields, values, err := parseProperties(model, r.idFieldName)
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

func (r SQLRepository[T]) Read(id any) (*T, error) {
	var result T
	fields, err := parseFieldNames(reflect.TypeOf(result))
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

	err = stmt.QueryRowx(id).StructScan(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r SQLRepository[T]) ReadAll() ([]T, error) {
	var result []T
	fields, err := parseFieldNames(reflect.TypeOf(result).Elem())
	if err != nil {
		return nil, err
	}

	sql := getSelectAllStmt(r.tableName, fields...)
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}

	err = stmt.Select(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r SQLRepository[T]) Update(id any, model T) error {
	fields, values, err := parseProperties(model, r.idFieldName)
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

func (r SQLRepository[T]) Delete(id any) error {
	sql, err := getDeleteFromStmt(r.dialect, r.tableName, r.idFieldName)
	if err != nil {
		return err
	}

	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return err
	}

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
