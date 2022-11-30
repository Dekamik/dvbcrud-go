package crudsql

import (
	"fmt"
	"strings"
)

type sqlTemplates interface {
	GetSelect() string
	GetSelectAll() string
	GetInsert(fields []string) (string, error)
	GetUpdate(fields []string) (string, error)
	GetDelete() string
}

type sqlTemplatesImpl struct {
	gen       paramGen
	tableName string
	idField   string

	selectSql    string
	selectAllSql string
	deleteSql    string
}

// GetSelect returns the SELECT statement (WHERE ID)
func (s sqlTemplatesImpl) GetSelect() string {
	return s.selectSql
}

// GetSelectAll returns the SELECT statement (all rows)
func (s sqlTemplatesImpl) GetSelectAll() string {
	return s.selectAllSql
}

// GetInsert generates and returns an INSERT INTO statement
func (s sqlTemplatesImpl) GetInsert(fields []string) (string, error) {
	placeholders, err := s.gen.GetParamPlaceholders(len(fields), Values)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		s.tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", ")), nil
}

// GetUpdate generates and returns an UPDATE statement
func (s sqlTemplatesImpl) GetUpdate(fields []string) (string, error) {
	columnPlaceholders, err := s.gen.GetParamPlaceholders(1, Columns)
	if err != nil {
		return "", err
	}

	// Not handling this error because the same code is run above
	valuePlaceholders, err := s.gen.GetParamPlaceholders(len(fields), Values)
	if err != nil {
		return "", err
	}

	f := fields
	for i := range f {
		f[i] += " = " + valuePlaceholders[i]
	}

	return fmt.Sprintf("UPDATE %s SET (%s) WHERE %s = %s",
		s.tableName,
		strings.Join(f, ", "),
		s.idField,
		columnPlaceholders[0]), nil
}

// GetDelete returns the DELETE statement
func (s sqlTemplatesImpl) GetDelete() string {
	return s.deleteSql
}

// generateSelect generates and returns a SELECT statement (WHERE ID)
func generateSelect(gen paramGen, tableName string, idField string, fields []string) (string, error) {
	placeholders, err := gen.GetParamPlaceholders(1, Columns)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s",
		strings.Join(fields, ", "),
		tableName,
		idField,
		placeholders[0]), nil
}

// generateSelectAll generates and returns a SELECT statement (all rows)
func generateSelectAll(table string, fields []string) string {
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), table)
}

// generateDelete returns DELETE FROM <table> WHERE <id> = ?
func generateDelete(gen paramGen, table string, idField string) (string, error) {
	placeholder, err := gen.GetParamPlaceholders(1, Columns)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s = %s",
		table,
		idField,
		placeholder[0]), nil
}

// newSQLTemplates pre-generates the SELECT, SELECT ALL and DELETE statement and returns a struct containing the templates.
func newSQLTemplates(gen paramGen, tableName string, idField string, allFields []string) (sqlTemplates, error) {
	selectSql, err := generateSelect(gen, tableName, idField, allFields)
	if err != nil {
		return nil, err
	}

	selectAllSql := generateSelectAll(tableName, allFields)

	deleteSql, err := generateDelete(gen, tableName, idField)
	if err != nil {
		return nil, err
	}

	return &sqlTemplatesImpl{
		gen:       gen,
		tableName: tableName,
		idField:   idField,

		selectSql:    selectSql,
		selectAllSql: selectAllSql,
		deleteSql:    deleteSql,
	}, nil
}
