package dvbcrud

import (
	"reflect"
	"testing"
	"time"
)

func TestParseFields(t *testing.T) {
	expected := []string{"UserId", "Name", "Surname", "Birthdate", "CreatedAt"}
	actual, _ := parseFieldNames(reflect.TypeOf(testUser{}))

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual fields didn't match expected fields")
	}
}

func TestParseFieldsAndValues(t *testing.T) {
	user := testUser{
		Id:        1,
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
