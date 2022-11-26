package dvbcrud_go

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

// SqlRepository holds configuration data for the SqlRepository.
type SqlRepository[T any] struct {
	db          *sqlx.DB
	tableName   string
	idFieldName string
}

func (r SqlRepository[T]) Create(model T) error {
	fields, values := r.parseFieldsAndValues(model)
	sql := getInsertInto(r.tableName, fields...)

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
	sql := getSelectFrom(r.tableName, r.idFieldName)
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}

	var result T
	err = stmt.QueryRowx(id).StructScan(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r SqlRepository[T]) ReadAll() ([]T, error) {
	sql := getSelectFrom(r.tableName, "")
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}

	var result []T
	err = stmt.Select(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r SqlRepository[T]) Update(id any, model T) error {
	fields, values := r.parseFieldsAndValues(model)
	sql := getUpdate(r.tableName, r.idFieldName, fields...)

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
	sql := getDeleteFrom(r.tableName, r.idFieldName)
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

// parseFieldsAndValues reads the TModel and returns its fields and
// values as two slices. The slices are synchronized which means the
// field and its value share the same index in both slices.
func (r SqlRepository[T]) parseFieldsAndValues(model T) ([]string, []any) {
	val := reflect.ValueOf(&model).Elem()
	numField := val.NumField() - 1
	fields := make([]string, numField)
	values := make([]any, numField)

	offset := 0
	for i := 0; i < numField+1; i++ {
		field := val.Type().Field(i)
		var name string

		if tag := field.Tag.Get("db"); tag != "" {
			name = tag
		} else {
			name = field.Name
		}

		if name == r.idFieldName {
			offset--
			continue
		}
		fields[i+offset] = name
		values[i+offset] = val.Field(i).Interface()
	}

	return fields, values
}
