package dvbcrud

import (
	"fmt"
	"reflect"
	"testing"
)

func newSqlParameterGeneratorMock(err error) sqlParameterGeneratorMock {
	if err != nil {
		return sqlParameterGeneratorMock{
			GetParamPlaceholdersMock: func(amount int, typ parameterType) ([]string, error) {
				return nil, err
			},
		}
	}
	return sqlParameterGeneratorMock{
		GetParamPlaceholdersMock: func(amount int, typ parameterType) ([]string, error) {
			placeholders := make([]string, amount)
			for i := 0; i < amount; i++ {
				placeholders[i] = "?"
			}
			return placeholders, nil
		},
	}
}

func TestSqlGeneratorImpl_GenerateSelect(t *testing.T) {
	sqlParamGenMock := newSqlParameterGeneratorMock(nil)
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	expected := "SELECT col_1, col_2 FROM any_table WHERE id_col = ?"
	actual, _ := sqlGen.GenerateSelect("any_table", "id_col", []string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateSelect_GetParamPlaceholdersErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	sqlParamGenMock := newSqlParameterGeneratorMock(expected)
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	_, actual := sqlGen.GenerateSelect("any_table", "id_col", []string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateSelectAll(t *testing.T) {
	sqlGen := sqlGeneratorImpl{}

	expected := "SELECT col_1, col_2 FROM any_table"
	actual := sqlGen.GenerateSelectAll("any_table", []string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateInsert(t *testing.T) {
	sqlParamGenMock := newSqlParameterGeneratorMock(nil)
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	expected := "INSERT INTO any_table (col_1, col_2) VALUES (?, ?)"
	actual, _ := sqlGen.GenerateInsert("any_table", []string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateInsert_GetParamPlaceholdersErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	sqlParamGenMock := newSqlParameterGeneratorMock(expected)
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	_, actual := sqlGen.GenerateInsert("", []string{})

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateUpdate(t *testing.T) {
	sqlParamGenMock := newSqlParameterGeneratorMock(nil)
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	expected := "UPDATE any_table SET (col_1 = ?, col_2 = ?) WHERE id_col = ?"
	actual, _ := sqlGen.GenerateUpdate("any_table", "id_col", []string{"col_1", "col_2"})

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateUpdate_GetColumnsParamPlaceholdersErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	sqlParamGenMock := newSqlParameterGeneratorMock(expected)
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	_, actual := sqlGen.GenerateUpdate("", "", []string{})

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateUpdate_GetValuesParamPlaceholdersErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	sqlParamGenMock := sqlParameterGeneratorMock{
		GetParamPlaceholdersMock: func(amount int, typ parameterType) ([]string, error) {
			if typ == Values {
				return nil, expected
			}
			return []string{}, nil
		},
	}
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	_, actual := sqlGen.GenerateUpdate("", "", []string{})

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateDelete(t *testing.T) {
	sqlParamGenMock := newSqlParameterGeneratorMock(nil)
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	expected := "DELETE FROM any_table WHERE id_col = ?"
	actual, _ := sqlGen.GenerateDelete("any_table", "id_col")

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlGeneratorImpl_GenerateDelete_GetParamPlaceholdersErr(t *testing.T) {
	expected := fmt.Errorf("AnyError")
	sqlParamGenMock := newSqlParameterGeneratorMock(expected)
	sqlGen := sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}

	_, actual := sqlGen.GenerateDelete("", "")

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestNewSqlGenerator(t *testing.T) {
	sqlParamGenMock := sqlParameterGeneratorMock{}

	expected := &sqlGeneratorImpl{
		sqlParamGen: sqlParamGenMock,
	}
	actual := newSQLGenerator(sqlParamGenMock)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\nExpected %v\nbut got %v", expected, actual)
	}
}
