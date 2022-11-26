package dvbcrud_go

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

// SqlRepository holds configuration data for the SqlRepository.
type SqlRepository struct {
	db          *sqlx.DB
	tableName   string
	idFieldName string
}

func SqlCreate[TModel any](r *SqlRepository, model TModel) error {
	fields, values := parseFieldsAndValues(r, model)
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

func SqlRead[TModel any](r *SqlRepository, id any) (*TModel, error) {
	sql := getSelectFrom(r.tableName, r.idFieldName)
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}

	var result TModel
	err = stmt.QueryRowx(id).StructScan(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func SqlReadAll[TModel any](r *SqlRepository) ([]TModel, error) {
	sql := getSelectFrom(r.tableName, "")
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}

	var result []TModel
	err = stmt.Select(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SqlUpdate[TModel any](repository *SqlRepository, id any, model TModel) error {
	fields, values := parseFieldsAndValues(repository, model)
	sql := getUpdate(repository.tableName, repository.idFieldName, fields...)

	stmt, err := repository.db.Preparex(sql)
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

func SqlDelete(repository *SqlRepository, id any) error {
	sql := getDeleteFrom(repository.tableName, repository.idFieldName)
	stmt, err := repository.db.Preparex(sql)
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
func parseFieldsAndValues[TModel any](repository *SqlRepository, model TModel) ([]string, []any) {
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

		if name == repository.idFieldName {
			offset--
			continue
		}
		fields[i+offset] = name
		values[i+offset] = val.Field(i).Interface()
	}

	return fields, values
}
