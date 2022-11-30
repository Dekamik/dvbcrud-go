package dvbcrud

import (
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
	sqlGenMock := sqlGeneratorMock{
		generateSelectMock: func(table string, idField string, fields []string) (string, error) {
			return "AnySelect", nil
		},
		generateSelectAllMock: func(table string, fields []string) string {
			return "AnySelectAll"
		},
		generateDeleteMock: func(table string, idField string) (string, error) {
			return "AnyDelete", nil
		},
	}

	expected := sqlTemplatesImpl{
		sqlGen:       sqlGenMock,
		tableName:    "any_table",
		idField:      "id_col",
		selectSql:    "AnySelect",
		selectAllSql: "AnySelectAll",
		deleteSql:    "AnyDelete",
	}

	actual, _ := newSQLTemplates(sqlGenMock, "any_table", "id_col", []string{"id_col", "col_1", "col_2"})

	if expected.GetSelect() != actual.GetSelect() ||
		expected.GetSelectAll() != actual.GetSelectAll() ||
		expected.GetDelete() != actual.GetDelete() {
		t.Fatalf("\nExpected %v\nbut got %v", expected, actual)
	}
}
