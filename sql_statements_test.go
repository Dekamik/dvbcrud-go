package dvbcrud

import (
	"reflect"
	"testing"
)

func TestSelectAll(t *testing.T) {
	actual := getSelectAllStmt("users", "UserId", "Name", "Surname", "Birthdate", "CreatedAt")
	expected := "SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM users;"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}

func TestSelectWhere(t *testing.T) {
	actual := getSelectFromStmt("users", "UserId", "UserId", "Name", "Surname", "Birthdate", "CreatedAt")
	expected := "SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM users WHERE UserId = ?;"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestInsertInto(t *testing.T) {
	actual := getInsertIntoStmt("users", "Name", "Surname", "Birthdate", "CreatedAt")
	expected := "INSERT INTO users (Name, Surname, Birthdate, CreatedAt) VALUES (?, ?, ?, ?);"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}

func TestUpdate(t *testing.T) {
	actual := getUpdateStmt("users", "UserId", "Name", "Surname", "Birthdate", "CreatedAt")
	expected := "UPDATE users SET (Name = ?, Surname = ?, Birthdate = ?, CreatedAt = ?) WHERE UserId = ?;"
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
	_ = getUpdateStmt("users", "id", fieldsCopy...)
	if !reflect.DeepEqual(fields, fieldsCopy) {
		t.Fatalf("getUpdateStmt mustn't mutate fields array")
	}
}

func TestDelete(t *testing.T) {
	actual := getDeleteFromStmt("users", "UserId")
	expected := "DELETE FROM users WHERE UserId = ?;"
	if actual != expected {
		t.Fatalf("Expected \"%s\" but was \"%s\" instead", expected, actual)
	}
}
