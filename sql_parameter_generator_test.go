package dvbcrud

import (
	"reflect"
	"testing"
)

func TestSqlParameterGeneratorImpl_GetParamPlaceholders_MySQL(t *testing.T) {
	gen := newSQLParamGen(MySQL)
	expected := []string{"?", "?", "?"}

	actual, _ := gen.GetParamPlaceholders(3, Columns)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlParameterGeneratorImpl_GetParamPlaceholders_PostgreSQL(t *testing.T) {
	gen := newSQLParamGen(PostgreSQL)
	expected := []string{"$1", "$2", "$3"}

	actual, _ := gen.GetParamPlaceholders(3, Columns)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlParameterGeneratorImpl_GetParamPlaceholders_OracleCol(t *testing.T) {
	gen := newSQLParamGen(Oracle)
	expected := []string{":col1", ":col2", ":col3"}

	actual, _ := gen.GetParamPlaceholders(3, Columns)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlParameterGeneratorImpl_GetParamPlaceholders_OracleVal(t *testing.T) {
	gen := newSQLParamGen(Oracle)
	expected := []string{":val1", ":val2", ":val3"}

	actual, _ := gen.GetParamPlaceholders(3, Values)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlParameterGeneratorImpl_GetParamPlaceholders_OracleColSingle(t *testing.T) {
	gen := newSQLParamGen(Oracle)
	expected := []string{":col"}

	actual, _ := gen.GetParamPlaceholders(1, Columns)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlParameterGeneratorImpl_GetParamPlaceholders_OracleValSingle(t *testing.T) {
	gen := newSQLParamGen(Oracle)
	expected := []string{":val"}

	actual, _ := gen.GetParamPlaceholders(1, Values)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestSqlParameterGeneratorImpl_GetParamPlaceholders_UnknownDialect(t *testing.T) {
	gen := newSQLParamGen(-1)
	expected := "unknown dialect"

	_, actual := gen.GetParamPlaceholders(1, Columns)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}
