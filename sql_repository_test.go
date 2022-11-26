package dvbcrud_go

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
	"time"
)

type testUser struct {
	Id        int `db:"user_id"`
	Name      string
	Surname   string
	Birthdate time.Time
	CreatedAt time.Time `db:"created_at"`
}

func newMock() (SqlRepository, *sql.DB, sqlmock.Sqlmock, error) {
	mockDb, mock, err := sqlmock.New()
	sqlxDb := sqlx.NewDb(mockDb, "sqlmock")
	repository := SqlRepository{
		db:          sqlxDb,
		tableName:   "users",
		idFieldName: "user_id",
	}
	return repository, mockDb, mock, err
}

func TestSqlRepository_Read(t *testing.T) {
	repo, mockDb, mock, _ := newMock()
	defer mockDb.Close()

	expected := testUser{
		Id:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"user_id", "name", "surname", "birthdate", "created_at"}).
		AddRow(expected.Id, expected.Name, expected.Surname, expected.Birthdate, expected.CreatedAt)
	mock.ExpectPrepare("^SELECT \\* FROM users WHERE user_id = \\?;$").
		ExpectQuery().
		WithArgs(expected.Id).
		WillReturnRows(rows)

	actual, err := Read[testUser](&repo, 1)
	if err != nil {
		t.Fatalf("Error on read: %s", err)
	}

	if !reflect.DeepEqual(&expected, actual) {
		t.Fatalf("Actual user didn't match actual user on read")
	}
}
