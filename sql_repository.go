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

func Read[TModel any](r *SqlRepository, id any) (*TModel, error) {
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

func ReadAll[TModel any](r *SqlRepository) ([]TModel, error) {
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

func Update[TModel any](repository *SqlRepository, id any, model TModel) error {
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

func Delete(repository *SqlRepository, id any) error {
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
	val := reflect.ValueOf(model).Elem()
	numField := val.NumField()
	fields := make([]string, numField)
	values := make([]any, numField)

	for i := 0; i < numField; i++ {
		name := val.Type().Field(i).Name
		if name == repository.idFieldName {
			continue
		}
		fields[i] = name
		values[i] = val.Field(i)
	}

	return fields, values
}
