package dvbcrud

import (
	"fmt"
	"testing"
)

func TestSqlTemplatesImpl_GetSelect(t *testing.T) {
	expected := "AnySelectStatement"
	templates := sqlTemplatesImpl{
		selectSql: expected,
	}

	actual := templates.GetSelect()

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetSelectAll(t *testing.T) {
	expected := "AnySelectAllStatement"
	templates := sqlTemplatesImpl{
		selectAllSql: expected,
	}

	actual := templates.GetSelectAll()

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetInsert(t *testing.T) {
	expected := "AnyInsertStatement"
	sqlGenMock := sqlGeneratorMock{
		generateInsertMock: func(table string, fields []string) (string, error) {
			return expected, nil
		},
	}
	templates := sqlTemplatesImpl{
		sqlGen: sqlGenMock,
	}

	actual, _ := templates.GetInsert([]string{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetUpdate(t *testing.T) {
	expected := "AnyUpdateStatement"
	sqlGenMock := sqlGeneratorMock{
		generateUpdateMock: func(table string, idField string, fields []string) (string, error) {
			return expected, nil
		},
	}
	templates := sqlTemplatesImpl{
		sqlGen: sqlGenMock,
	}

	actual, _ := templates.GetUpdate([]string{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetDelete(t *testing.T) {
	expected := "AnyDeleteStatement"
	templates := sqlTemplatesImpl{
		deleteSql: expected,
	}

	actual := templates.GetDelete()

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
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

func TestNewSQLTemplates_GenerateSelectErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	sqlGenMock := sqlGeneratorMock{
		generateSelectMock: func(table string, idField string, fields []string) (string, error) {
			return "", expected
		},
	}

	_, actual := newSQLTemplates(sqlGenMock, "", "", []string{})

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestNewSQLTemplates_GenerateDeleteErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	sqlGenMock := sqlGeneratorMock{
		generateSelectMock: func(table string, idField string, fields []string) (string, error) {
			return "", nil
		},
		generateSelectAllMock: func(table string, fields []string) string {
			return ""
		},
		generateDeleteMock: func(table string, idField string) (string, error) {
			return "", expected
		},
	}

	_, actual := newSQLTemplates(sqlGenMock, "", "", []string{})

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}
