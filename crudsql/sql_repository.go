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
	templates   sqlTemplates
	idFieldName string
}

// Create inserts the values in model into a new row in the table.
func (r SQLRepository[T]) Create(model T) error {
	fields, values, err := internal.ParseProperties(model, r.idFieldName)
	if err != nil {
		return err
	}

	sql, err := r.templates.GetInsert(fields)
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
	sql := r.templates.GetSelect()
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var result T
	err = stmt.QueryRowx(id).StructScan(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ReadAll fetches all rows from the table.
func (r SQLRepository[T]) ReadAll() ([]T, error) {
	sql := r.templates.GetSelectAll()
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var result []T
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

	sql, err := r.templates.GetUpdate(fields)
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
	sql := r.templates.GetDelete()
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

// New creates and returns a new SQLRepository.
func New[T any](db *sqlx.DB, dialect SQLDialect, tableName string, idFieldName string) (*SQLRepository[T], error) {
	if db == nil {
		return nil, fmt.Errorf("db cannot be nil")
	}
	if tableName == "" {
		return nil, fmt.Errorf("tableName cannot be empty")
	}
	if idFieldName == "" {
		idFieldName = "id"
	}

	var hack T
	fields, err := internal.ParseFieldNames(reflect.TypeOf(hack))
	if err != nil {
		return nil, err
	}

	paramGen := NewSQLParamGen(dialect)
	statementGen, err := newSQLTemplates(paramGen, tableName, idFieldName, fields)

	return &SQLRepository[T]{
		db:          db,
		templates:   statementGen,
		idFieldName: idFieldName,
	}, nil
}
