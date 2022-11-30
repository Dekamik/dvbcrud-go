package dvbcrud

import (
	"fmt"
	"testing"
)

func TestSqlTemplatesImpl_GetInsert(t *testing.T) {
	sqlGenMock := sqlGeneratorMock{
		generateInsertMock: func(table string, fields []string) (string, error) {
			return "AnyInsert", nil
		},
	}
	mock := sqlTemplatesImpl{
		sqlGen:    sqlGenMock,
		tableName: "any_table",
	}
	expected := "AnyInsert"

	actual, _ := mock.GetInsert([]string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetInsertParamsErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	sqlGenMock := sqlGeneratorMock{
		generateInsertMock: func(table string, fields []string) (string, error) {
			return "", expected
		},
	}
	mock := sqlTemplatesImpl{sqlGen: sqlGenMock}

	_, actual := mock.GetInsert([]string{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\"", expected, actual)
	}
}

func TestSqlTemplatesImpl_GetUpdate(t *testing.T) {
	sqlGenMock := sqlGeneratorMock{
		generateUpdateMock: func(table string, idField string, fields []string) (string, error) {
			return "AnyUpdate", nil
		},
	}
	mock := sqlTemplatesImpl{
		sqlGen:    sqlGenMock,
		tableName: "any_table",
		idField:   "id_col",
	}
	expected := "AnyUpdate"

	actual, _ := mock.GetUpdate([]string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}
