package dvbcrud

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
	"time"
)

type repoTestUser struct {
	Id        uint64    `db:"UserId"`
	Name      string    `db:"Name"`
	Surname   string    `db:"Surname"`
	Birthdate time.Time `db:"Birthdate"`
	CreatedAt time.Time `db:"CreatedAt"`
}

func newMock() (*SqlRepository[repoTestUser], *sql.DB, sqlmock.Sqlmock, error) {
	mockDb, mock, err := sqlmock.New()
	sqlxDb := sqlx.NewDb(mockDb, "sqlmock")
	repo, _ := NewSql[repoTestUser](sqlxDb, "Users", "UserId")
	return repo, mockDb, mock, err
}

func TestSqlRepository_Create(t *testing.T) {
	repo, mockDb, mock, _ := newMock()
	defer mockDb.Close()

	user := repoTestUser{
		Id:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	mock.ExpectPrepare("^INSERT INTO Users \\(Name, Surname, Birthdate, CreatedAt\\) VALUES \\(\\?, \\?, \\?, \\?\\);$").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Expected Create to succeed, but got: %s", err)
	}
}

func TestSqlRepository_Read(t *testing.T) {
	repo, mockDb, mock, _ := newMock()
	defer mockDb.Close()

	expected := repoTestUser{
		Id:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"UserId", "Name", "Surname", "Birthdate", "CreatedAt"}).
		AddRow(expected.Id, expected.Name, expected.Surname, expected.Birthdate, expected.CreatedAt)
	mock.ExpectPrepare("^SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM Users WHERE UserId = \\?;$").
		ExpectQuery().
		WithArgs(expected.Id).
		WillReturnRows(rows)

	actual, err := repo.Read(1)
	if err != nil {
		t.Fatalf("Error on Read: %s", err)
	}

	if !reflect.DeepEqual(&expected, actual) {
		t.Fatalf("Actual user must match expected user on Read")
	}
}

func TestSqlRepository_ReadAll(t *testing.T) {
	repo, mockDb, mock, _ := newMock()
	defer mockDb.Close()

	expected := []repoTestUser{
		{
			Id:        1,
			Name:      "AnyName1",
			Surname:   "AnySurname1",
			Birthdate: time.Now(),
			CreatedAt: time.Now(),
		},
		{
			Id:        2,
			Name:      "AnyName2",
			Surname:   "AnySurname2",
			Birthdate: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"UserId", "Name", "Surname", "Birthdate", "CreatedAt"}).
		AddRow(expected[0].Id, expected[0].Name, expected[0].Surname, expected[0].Birthdate, expected[0].CreatedAt).
		AddRow(expected[1].Id, expected[1].Name, expected[1].Surname, expected[1].Birthdate, expected[1].CreatedAt)
	mock.ExpectPrepare("^SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM Users;$").
		ExpectQuery().
		WillReturnRows(rows)

	actual, err := repo.ReadAll()
	if err != nil {
		t.Fatalf("Error on ReadAll: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual users must match expected users on ReadAll")
	}
}

func TestSqlRepository_Update(t *testing.T) {
	repo, mockDb, mock, _ := newMock()
	defer mockDb.Close()

	user := repoTestUser{
		Id:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	mock.ExpectPrepare("^UPDATE Users SET \\(Name = \\?, Surname = \\?, Birthdate = \\?, CreatedAt = \\?\\) WHERE UserId = \\?;$").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(1, user)
	if err != nil {
		t.Fatalf("Expected Update to succeed, but got: %s", err)
	}
}

func TestSqlRepository_Delete(t *testing.T) {
	repo, mockDb, mock, _ := newMock()
	defer mockDb.Close()

	mock.ExpectPrepare("^DELETE FROM Users WHERE UserId = \\?;$").
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(1)
	if err != nil {
		t.Fatalf("Expected Delete to succeed, but got: %s", err)
	}
}

func TestNewSql(t *testing.T) {
	mockDb, _, _ := sqlmock.New()
	sqlxDb := sqlx.NewDb(mockDb, "sqlmock")
	repo, _ := NewSql[repoTestUser](sqlxDb, "Users", "UserId")

	if repo == nil {
		t.Fatalf("Expected a repo, but got nil instead")
	}
}

func TestNewSqlNilDb(t *testing.T) {
	_, err := NewSql[repoTestUser](nil, "users", "UserId")
	if err == nil {
		t.Fatalf("Expected error on nil db")
	}

	expected := "db cannot be nil"
	if err.Error() != expected {
		t.Fatalf("Expected \"%s\" error but got \"%s\" instead", expected, err.Error())
	}
}

func TestNewSqlEmptyTableName(t *testing.T) {
	mockDb, _, _ := sqlmock.New()
	sqlxDb := sqlx.NewDb(mockDb, "sqlmock")
	_, err := NewSql[repoTestUser](sqlxDb, "", "UserId")
	if err == nil {
		t.Fatalf("Expected error on empty table name")
	}

	expected := "tableName cannot be empty"
	if err.Error() != expected {
		t.Fatalf("Expected \"%s\" error but got \"%s\" instead", expected, err.Error())
	}
}

func TestNewSqlEmptyIdFieldName(t *testing.T) {
	mockDb, _, _ := sqlmock.New()
	sqlxDb := sqlx.NewDb(mockDb, "sqlmock")
	repo, _ := NewSql[repoTestUser](sqlxDb, "users", "")
	if repo == nil {
		t.Fatalf("Expected a repo on empty idFieldName")
	}

	if repo.idFieldName != "id" {
		t.Fatalf("Expected idFieldName to be \"id\"")
	}
}
