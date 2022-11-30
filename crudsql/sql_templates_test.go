package crudsql

import (
	"reflect"
	"testing"
)

type paramGenMock struct {
	paramGen
	OnReturn []string
	OnError  error
}

func (p paramGenMock) GetParamPlaceholders(amount int, typ paramType) ([]string, error) {
	if p.OnError != nil {
		return nil, p.OnError
	}
	return p.OnReturn, nil
}

func TestSqlTemplatesImpl_GetSelect(t *testing.T) {
	expected := "AnySelectStatement"
	mock := sqlTemplatesImpl{
		selectSql: expected,
	}

	actual := mock.GetSelect()

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetSelectAll(t *testing.T) {
	expected := "AnySelectAllStatement"
	mock := sqlTemplatesImpl{
		selectAllSql: expected,
	}

	actual := mock.GetSelectAll()

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetInsert(t *testing.T) {
	paramsGenMock := paramGenMock{
		OnReturn: []string{"?", "?"},
	}
	mock := sqlTemplatesImpl{
		gen:       paramsGenMock,
		tableName: "any_table",
	}
	expected := "INSERT INTO any_table (col_1, col_2) VALUES (?, ?)"

	actual, _ := mock.GetInsert([]string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetUpdate(t *testing.T) {
	paramsGenMock := paramGenMock{
		OnReturn: []string{"?", "?"},
	}
	mock := sqlTemplatesImpl{
		gen:       paramsGenMock,
		tableName: "any_table",
		idField:   "id_col",
	}
	expected := "UPDATE any_table SET (col_1 = ?, col_2 = ?) WHERE id_col = ?"

	actual, _ := mock.GetUpdate([]string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetDelete(t *testing.T) {
	expected := "AnyDeleteStatement"
	mock := sqlTemplatesImpl{
		deleteSql: expected,
	}

	actual := mock.GetDelete()

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestNewSQLTemplates(t *testing.T) {
	mock := paramGenMock{
		OnReturn: []string{"?"},
	}

	expected := &sqlTemplatesImpl{
		gen:          mock,
		tableName:    "any_table",
		idField:      "id_col",
		selectSql:    "SELECT id_col, col_1, col_2 FROM any_table WHERE id_col = ?",
		selectAllSql: "SELECT id_col, col_1, col_2 FROM any_table",
		deleteSql:    "DELETE FROM any_table WHERE id_col = ?",
	}

	actual, _ := newSQLTemplates(mock, "any_table", "id_col", []string{"id_col", "col_1", "col_2"})

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected %v\nbut got %v", expected, actual)
	}
}
