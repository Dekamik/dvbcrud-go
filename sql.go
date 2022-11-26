package dvbcrud_go

import (
    "errors"
    "fmt"
    "github.com/jmoiron/sqlx"
    "reflect"
    "strings"
)

type SqlConfig struct {
    db        *sqlx.DB
    tableName string
    idName    string
}

type SqlRepository[TModel any] struct {
    config SqlConfig
}

func (r SqlRepository[TModel]) Create(model TModel) error {
    fields, values := r.parseFieldsAndValues(model)

    sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
        r.config.tableName,
        strings.Join(fields, ", "),
        strings.Repeat("?, ", len(values)-1)+"?")
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

func (r SqlRepository[TModel]) Read(id any) (TModel, error) {
    sql := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?",
        r.config.tableName,
        r.config.idName)
    stmt, err := r.config.db.Preparex(sql)
    if err != nil {
        return nil, err
    }

    var result TModel
    err = stmt.Select(&result, id)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (r SqlRepository[TModel]) ReadAll() ([]TModel, error) {
    sql := fmt.Sprintf("SELECT * FROM %s", r.config.tableName)
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
    return nil
}

func (r SqlRepository[TModel]) Delete(id any) error {
    sql := fmt.Sprintf("DELETE * FROM %s WHERE %s = ?",
        r.config.tableName,
        r.config.idName)
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
        return errors.New("more than 1 row affected by delete statement")
    }

    return nil
}

func NewSql[TModel any](config SqlConfig) (Repository[TModel], error) {
    if config.db == nil {
        return nil, errors.New("config.db must be set")
    }
    if config.tableName == "" {
        return nil, errors.New("config.tableName must be defined")
    }
    if config.idName == "" {
        config.idName = "id"
    }

    return SqlRepository[TModel]{
        config: config,
    }, nil
}

func (r SqlRepository[TModel]) parseFieldsAndValues(model TModel) ([]string, []any) {
    val := reflect.ValueOf(model).Elem()
    numField := val.NumField()
    fields := make([]string, numField)
    values := make([]any, numField)

    for i := 0; i < numField; i++ {
        name := val.Type().Field(i).Name
        if name == r.config.idName {
            continue
        }
        fields[i] = name
        values[i] = val.Field(i)
    }

    return fields, values
}
