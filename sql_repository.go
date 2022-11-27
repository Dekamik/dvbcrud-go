package dvbcrud

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

// SqlRepository handles queries to a table in an SQL database.
// T is the struct type that will be mapped against the table rows.
type SqlRepository[T any] struct {
	db          *sqlx.DB
	tableName   string
	idFieldName string
}

func (r SqlRepository[T]) Create(model T) error {
	fields, values := parseProperties(model, r.idFieldName)
	sql := getInsertIntoStmt(r.tableName, fields...)

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

func (r SqlRepository[T]) Read(id any) (*T, error) {
	var result T
	fields := parseFieldNames(reflect.TypeOf(result))
	sql := getSelectFromStmt(r.tableName, r.idFieldName, fields...)
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

func (r SqlRepository[T]) ReadAll() ([]T, error) {
	var result []T
	fields := parseFieldNames(reflect.TypeOf(result).Elem())
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

func (r SqlRepository[T]) Update(id any, model T) error {
	fields, values := parseProperties(model, r.idFieldName)
	sql := getUpdateStmt(r.tableName, r.idFieldName, fields...)

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

func (r SqlRepository[T]) Delete(id any) error {
	sql := getDeleteFromStmt(r.tableName, r.idFieldName)
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
		return errors.New("more than 1 row affected by DELETE statement")
	}

	return nil
}

func NewSql[T any](db *sqlx.DB, tableName string, idFieldName string) (*SqlRepository[T], error) {
	if db == nil {
		return nil, fmt.Errorf("db cannot be nil")
	}
	if tableName == "" {
		return nil, fmt.Errorf("tableName cannot be empty")
	}
	if idFieldName == "" {
		idFieldName = "id"
	}

	return &SqlRepository[T]{
		db:          db,
		tableName:   tableName,
		idFieldName: idFieldName,
	}, nil
}

func parseFieldNames[T reflect.Type](typ T) []string {
	numFields := typ.NumField()
	fields := make([]string, numFields)

	for i := 0; i < numFields; i++ {
		fields[i] = typ.Field(i).Tag.Get("db")
	}

	return fields
}

// parseProperties reads the struct type T and returns its fields
// and values as two slices. The slices are synchronized which means each
// field and its corresponding value share the same index in both slices.
// Specifying idFieldName filters out that field, which is useful in
// INSERTS and UPDATES.
func parseProperties[T any](model T, idFieldName string) ([]string, []any) {
	val := reflect.ValueOf(&model).Elem()
	numField := val.NumField()
	fields := make([]string, numField-1)
	values := make([]any, numField-1)

	index := 0
	for i := 0; i < numField; i++ {
		name := val.Type().Field(i).Tag.Get("db")
		if name == idFieldName {
			continue
		}

		fields[index] = name
		values[index] = val.Field(i).Interface()
		index++
	}

	return fields, values
}
