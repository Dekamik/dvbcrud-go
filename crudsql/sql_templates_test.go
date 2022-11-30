package crudsql

import (
	"fmt"
	"testing"
)

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
		GetParamPlaceholdersMock: func(amount int, typ paramType) ([]string, error) {
			return []string{"?", "?"}, nil
		},
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

func TestSqlTemplatesImpl_GetInsertParamsErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	paramGenMock := paramGenMock{
		GetParamPlaceholdersMock: func(amount int, typ paramType) ([]string, error) {
			return nil, expected
		},
	}
	mock := sqlTemplatesImpl{gen: paramGenMock}

	_, actual := mock.GetInsert([]string{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetUpdate(t *testing.T) {
	paramsGenMock := paramGenMock{
		GetParamPlaceholdersMock: func(amount int, typ paramType) ([]string, error) {
			return []string{"?", "?"}, nil
		},
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
		GetParamPlaceholdersMock: func(amount int, typ paramType) ([]string, error) {
			return []string{"?"}, nil
		},
	}

	expected := sqlTemplatesImpl{
		gen:          mock,
		tableName:    "any_table",
		idField:      "id_col",
		selectSql:    "SELECT id_col, col_1, col_2 FROM any_table WHERE id_col = ?",
		selectAllSql: "SELECT id_col, col_1, col_2 FROM any_table",
		deleteSql:    "DELETE FROM any_table WHERE id_col = ?",
	}

	actual, _ := newSQLTemplates(mock, "any_table", "id_col", []string{"id_col", "col_1", "col_2"})

	if expected.GetSelect() != actual.GetSelect() ||
		expected.GetSelectAll() != actual.GetSelectAll() ||
		expected.GetDelete() != actual.GetDelete() {
		t.Fatalf("\nExpected %v\nbut got %v", expected, actual)
	}
}
