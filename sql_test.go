package dvbcrud_go

import (
	"reflect"
	"testing"
)

func TestSelectAll(t *testing.T) {
	actual := getSelectFrom("users", "")
	expected := "SELECT * FROM users;"
	if actual != expected {
		t.Fatalf("getSelectFrom(\"users\", \"\") should be \"%s\" but was %s instead", expected, actual)
	}
}

func TestSelectWhere(t *testing.T) {
	actual := getSelectFrom("users", "id")
	expected := "SELECT * FROM users WHERE id = ?;"
	if actual != expected {
		t.Fatalf("getSelectFrom(\"users\", \"id\") should be \"%s\" but was %s instead", expected, actual)
	}
}

func TestInsertInto(t *testing.T) {
	actual := getInsertInto("users", "name", "surname", "birthdate", "created_at")
	expected := "INSERT INTO users (name, surname, birthdate, created_at) VALUES (?, ?, ?, ?);"
	if actual != expected {
		t.Fatalf("getInsertInto(\"users\", \"name\", \"surname\", \"birthdate\", \"created_at\") should be \"%s\" but was %s instead", expected, actual)
	}
}

func TestUpdate(t *testing.T) {
	actual := getUpdate("users", "id", "name", "surname", "birthdate", "created_at")
	expected := "UPDATE users SET (name = ?, surname = ?, birthdate = ?, created_at = ?) WHERE id = ?;"
	if actual != expected {
		t.Fatalf("getUpdate(\"users\", \"id\", \"name\", \"surname\", \"birthdate\", \"created_at\") should be \"%s\" but was %s instead", expected, actual)
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
	_ = getUpdate("users", "id", fieldsCopy...)
	if !reflect.DeepEqual(fields, fieldsCopy) {
		t.Fatalf("getUpdate mustn't mutate fields array")
	}
}

func TestDelete(t *testing.T) {
	actual := getDeleteFrom("users", "id")
	expected := "DELETE * FROM users WHERE id = ?;"
	if actual != expected {
		t.Fatalf("getDeleteFrom(\"users\", \"id\") should be \"%s\" but was %s instead", expected, actual)
	}
}
