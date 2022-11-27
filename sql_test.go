package dvbcrud_go

import (
	"reflect"
	"testing"
)

func TestSelectAll(t *testing.T) {
	actual := getSelectFromStmt("users", "")
	expected := "SELECT * FROM users;"
	if actual != expected {
		t.Fatalf("getSelectFromStmt(\"users\", \"\") should be \"%s\" but was %s instead", expected, actual)
	}
}

func TestSelectWhere(t *testing.T) {
	actual := getSelectFromStmt("users", "id")
	expected := "SELECT * FROM users WHERE id = ?;"
	if actual != expected {
		t.Fatalf("getSelectFromStmt(\"users\", \"id\") should be \"%s\" but was %s instead", expected, actual)
	}
}

func TestInsertInto(t *testing.T) {
	actual := getInsertIntoStmt("users", "name", "surname", "birthdate", "created_at")
	expected := "INSERT INTO users (name, surname, birthdate, created_at) VALUES (?, ?, ?, ?);"
	if actual != expected {
		t.Fatalf("getInsertIntoStmt(\"users\", \"name\", \"surname\", \"birthdate\", \"created_at\") should be \"%s\" but was %s instead", expected, actual)
	}
}

func TestUpdate(t *testing.T) {
	actual := getUpdateStmt("users", "id", "name", "surname", "birthdate", "created_at")
	expected := "UPDATE users SET (name = ?, surname = ?, birthdate = ?, created_at = ?) WHERE id = ?;"
	if actual != expected {
		t.Fatalf("getUpdateStmt(\"users\", \"id\", \"name\", \"surname\", \"birthdate\", \"created_at\") should be \"%s\" but was %s instead", expected, actual)
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
	_ = getUpdateStmt("users", "id", fieldsCopy...)
	if !reflect.DeepEqual(fields, fieldsCopy) {
		t.Fatalf("getUpdateStmt mustn't mutate fields array")
	}
}

func TestDelete(t *testing.T) {
	actual := getDeleteFromStmt("users", "id")
	expected := "DELETE * FROM users WHERE id = ?;"
	if actual != expected {
		t.Fatalf("getDeleteFromStmt(\"users\", \"id\") should be \"%s\" but was %s instead", expected, actual)
	}
}
