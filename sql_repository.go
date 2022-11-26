package dvbcrud_go

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

// SqlRepositoryConfig holds configuration data for the SqlRepository.
type SqlRepositoryConfig struct {
	db          *sqlx.DB
	tableName   string
	idFieldName string
}

type SqlRepository[TModel any] struct {
	config SqlRepositoryConfig
	Repository[TModel]
}

func (r SqlRepository[TModel]) Create(model TModel) error {
	fields, values := r.parseFieldsAndValues(model)
	sql := getInsertInto(r.config.tableName, fields...)

	stmt, err := r.config.db.Preparex(sql)
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

func (r SqlRepository[TModel]) Read(id any) (*TModel, error) {
	sql := getSelectFrom(r.config.tableName, r.config.idFieldName)
	stmt, err := r.config.db.Preparex(sql)
	if err != nil {
		return nil, err
	}

	var result TModel
	err = stmt.Get(&result, id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r SqlRepository[TModel]) ReadAll() ([]TModel, error) {
	sql := getSelectFrom(r.config.tableName, "")
	stmt, err := r.config.db.Preparex(sql)
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

func (r SqlRepository[TModel]) Update(id any, model TModel) error {
	fields, values := r.parseFieldsAndValues(model)
	sql := getUpdate(r.config.tableName, r.config.idFieldName, fields...)

	stmt, err := r.config.db.Preparex(sql)
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

func (r SqlRepository[TModel]) Delete(id any) error {
	sql := getDeleteFrom(r.config.tableName, r.config.idFieldName)
	stmt, err := r.config.db.Preparex(sql)
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

// NewSql returns a new SqlRepository struct based on the configuration
// in the SqlRepositoryConfig.
func NewSql[TModel any](config SqlRepositoryConfig) (Repository[TModel], error) {
	if config.db == nil {
		return nil, errors.New("config.db must be set")
	}
	if config.tableName == "" {
		return nil, errors.New("config.tableName must be defined")
	}
	if config.idFieldName == "" {
		config.idFieldName = "id"
	}

	return SqlRepository[TModel]{
		config: config,
	}, nil
}

// parseFieldsAndValues reads the TModel and returns its fields and
// values as two slices. The slices are synchronized which means the
// field and its value share the same index in both slices.
func (r SqlRepository[TModel]) parseFieldsAndValues(model TModel) ([]string, []any) {
	val := reflect.ValueOf(model).Elem()
	numField := val.NumField()
	fields := make([]string, numField)
	values := make([]any, numField)

	for i := 0; i < numField; i++ {
		name := val.Type().Field(i).Name
		if name == r.config.idFieldName {
			continue
		}
		fields[i] = name
		values[i] = val.Field(i)
	}

	return fields, values
}
