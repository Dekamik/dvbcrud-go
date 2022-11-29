package crudsql

import (
	"reflect"
	"testing"
)

func TestGetParamPlaceholdersMySQL(t *testing.T) {
	expected := []string{"?", "?", "?"}
	actual, _ := getParamPlaceholders(MySQL, Columns, 3)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestGetParamPlaceholdersPostgreSQL(t *testing.T) {
	expected := []string{"$1", "$2", "$3"}
	actual, _ := getParamPlaceholders(PostgreSQL, Columns, 3)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestGetParamPlaceholdersOracleCol(t *testing.T) {
	expected := []string{":col1", ":col2", ":col3"}
	actual, _ := getParamPlaceholders(Oracle, Columns, 3)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestGetParamPlaceholdersOracleVal(t *testing.T) {
	expected := []string{":val1", ":val2", ":val3"}
	actual, _ := getParamPlaceholders(Oracle, Values, 3)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestGetParamPlaceholdersOracleColSingle(t *testing.T) {
	expected := []string{":col"}
	actual, _ := getParamPlaceholders(Oracle, Columns, 1)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestGetParamPlaceholdersOracleValSingle(t *testing.T) {
	expected := []string{":val"}
	actual, _ := getParamPlaceholders(Oracle, Values, 1)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestGetParamPlaceholdersUnknownDialect(t *testing.T) {
	expected := "unknown dialect"
	_, actual := getParamPlaceholders(-1, Columns, 1)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestSelectAll(t *testing.T) {
	actual := getSelectAllStmt("users", "UserId", "Name", "Surname", "Birthdate", "CreatedAt")
	expected := "SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM users"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}

func TestSelectWhere(t *testing.T) {
	actual, _ := getSelectFromStmt(MySQL, "users", "UserId", "UserId", "Name", "Surname", "Birthdate", "CreatedAt")
	expected := "SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM users WHERE UserId = ?"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestInsertInto(t *testing.T) {
	actual, _ := getInsertIntoStmt(MySQL, "users", "Name", "Surname", "Birthdate", "CreatedAt")
	expected := "INSERT INTO users (Name, Surname, Birthdate, CreatedAt) VALUES (?, ?, ?, ?)"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}

func TestUpdate(t *testing.T) {
	actual, _ := getUpdateStmt(MySQL, "users", "UserId", "Name", "Surname", "Birthdate", "CreatedAt")
	expected := "UPDATE users SET (Name = ?, Surname = ?, Birthdate = ?, CreatedAt = ?) WHERE UserId = ?"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}

func TestUpdatePassByValue(t *testing.T) {
	fields := []string{
		"name",
		"surname",
		"birthdate",
		"created_at",
	}
	fieldsCopy := fields
	_, _ = getUpdateStmt(MySQL, "users", "id", fieldsCopy...)
	if !reflect.DeepEqual(fields, fieldsCopy) {
		t.Fatalf("getUpdateStmt mustn't mutate fields array")
	}
}

func TestDelete(t *testing.T) {
	actual, _ := getDeleteFromStmt(MySQL, "users", "UserId")
	expected := "DELETE FROM users WHERE UserId = ?"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}
