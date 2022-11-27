package dvbcrud

import (
	"reflect"
	"testing"
	"time"
)

type structTestUser struct {
	ID        uint64    `db:"UserId"`
	Name      string    `db:"Name"`
	Surname   string    `db:"Surname"`
	Birthdate time.Time `db:"Birthdate"`
	CreatedAt time.Time `db:"CreatedAt"`
}

type testMissingTagAddress struct {
	ID      uint64 `db:"address_id"`
	Address string `db:"address"`
	ZipCode string `db:"zip_code"`
	City    string
}

func TestParseFieldNames(t *testing.T) {
	expected := []string{"UserId", "Name", "Surname", "Birthdate", "CreatedAt"}
	actual, _ := parseFieldNames(reflect.TypeOf(structTestUser{}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual fields didn't match expected fields")
	}
}

func TestParseFieldNamesNonStructType(t *testing.T) {
	_, err := parseFieldNames(reflect.TypeOf([]int{}))
	if err == nil {
		t.Fatalf("Expected error on non-struct type")
	}
	expected := "type must be a kind of struct"
	if err.Error() != expected {
		t.Fatalf("Expected error \"%s\" but got \"%s\" instead", expected, err.Error())
	}
}

func TestParseFieldNamesMissingTag(t *testing.T) {
	_, err := parseFieldNames(reflect.TypeOf(testMissingTagAddress{}))
	if err == nil {
		t.Fatalf("Expected error on missing tag")
	}
	expected := "testMissingTagAddress.City lacks a db tag"
	if err.Error() != expected {
		t.Fatalf("Expected error \"%s\" but got \"%s\" instead", expected, err.Error())
	}
}

func TestParseProperties(t *testing.T) {
	user := structTestUser{
		ID:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	expectedFields := []string{"Name", "Surname", "Birthdate", "CreatedAt"}
	expectedValues := []any{user.Name, user.Surname, user.Birthdate, user.CreatedAt}
	actualFields, actualValues, _ := parseProperties(user, "UserId")

	if !reflect.DeepEqual(expectedFields, actualFields) {
		t.Fatalf("Actual fields didn't match expected fields")
	} else if !reflect.DeepEqual(expectedValues, actualValues) {
		t.Fatalf("Actual values didn't match expected values")
	}
}

func TestParsePropertiesNonStructType(t *testing.T) {
	test := []string{"one"}
	_, _, err := parseProperties(test, "")
	if err == nil {
		t.Fatalf("Expected error on non-struct type")
	}
	expected := "model must be a struct type"
	if err.Error() != expected {
		t.Fatalf("Expected error \"%s\" but got \"%s\" instead", expected, err.Error())
	}
}

func TestParsePropertiesMissingTag(t *testing.T) {
	address := testMissingTagAddress{
		ID:      0,
		Address: "",
		ZipCode: "",
		City:    "",
	}
	_, _, err := parseProperties(address, "address_id")
	if err == nil {
		t.Fatalf("Expected error on missing tag")
	}
	expected := "testMissingTagAddress.City lacks a db tag"
	if err.Error() != expected {
		t.Fatalf("Expected error \"%s\" but got \"%s\" instead", expected, err.Error())
	}
}
