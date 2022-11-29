package crudsql

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
	"time"
)

type repoTestUser struct {
	ID        uint64    `db:"UserId"`
	Name      string    `db:"Name"`
	Surname   string    `db:"Surname"`
	Birthdate time.Time `db:"Birthdate"`
	CreatedAt time.Time `db:"CreatedAt"`
}

type parseCrashingType struct {
	id any
}

func newMock[T any]() (*SQLRepository[T], *sql.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := NewSQLRepository[T](sqlxDb, MySQL, "Users", "UserId")
	return repo, mockDB, mock, err
}

func TestSqlRepository_Create(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()

	user := repoTestUser{
		ID:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	mock.ExpectPrepare("^INSERT INTO Users \\(Name, Surname, Birthdate, CreatedAt\\) VALUES \\(\\?, \\?, \\?, \\?\\)$").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Expected Create to succeed, but got: %s", err)
	}
}

func TestSqlRepository_CreateParsePropertiesErr(t *testing.T) {
	repo, mockDB, _, _ := newMock[parseCrashingType]()
	defer mockDB.Close()
	expected := "parseCrashingType.id lacks a db tag"

	model := parseCrashingType{}
	actual := repo.Create(model)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSQLRepository_CreateGetSqlErr(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := NewSQLRepository[repoTestUser](sqlxDb, -1, "Users", "UserId")

	expected := "unknown dialect"
	actual := repo.Create(repoTestUser{})

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_CreatePrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^INSERT INTO Users \\(Name, Surname, Birthdate, CreatedAt\\) VALUES \\(\\?, \\?, \\?, \\?\\)$").
		WillReturnError(expected)

	actual := repo.Create(repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_CreateExecErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^INSERT INTO Users \\(Name, Surname, Birthdate, CreatedAt\\) VALUES \\(\\?, \\?, \\?, \\?\\)$").
		ExpectExec().
		WillReturnError(expected)

	actual := repo.Create(repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_CreateRowsAffectedErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	user := repoTestUser{}
	mock.ExpectPrepare("^INSERT INTO Users \\(Name, Surname, Birthdate, CreatedAt\\) VALUES \\(\\?, \\?, \\?, \\?\\)$").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt).
		WillReturnResult(sqlmock.NewErrorResult(expected))

	actual := repo.Create(repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_CreateOtherThanOneRowAffected(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := "2 rows affected by INSERT INTO statement"
	user := repoTestUser{}
	mock.ExpectPrepare("^INSERT INTO Users \\(Name, Surname, Birthdate, CreatedAt\\) VALUES \\(\\?, \\?, \\?, \\?\\)$").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 2))

	actual := repo.Create(repoTestUser{})

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_Read(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()

	expected := repoTestUser{
		ID:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"UserId", "Name", "Surname", "Birthdate", "CreatedAt"}).
		AddRow(expected.ID, expected.Name, expected.Surname, expected.Birthdate, expected.CreatedAt)
	mock.ExpectPrepare("^SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM Users WHERE UserId = \\?$").
		ExpectQuery().
		WithArgs(expected.ID).
		WillReturnRows(rows)

	actual, err := repo.Read(1)
	if err != nil {
		t.Fatalf("Error on Read: %s", err)
	}

	if !reflect.DeepEqual(&expected, actual) {
		t.Fatalf("Actual user must match expected user on Read")
	}
}

func TestSqlRepository_ReadParsePropertiesErr(t *testing.T) {
	repo, mockDB, _, _ := newMock[parseCrashingType]()
	defer mockDB.Close()
	expected := "parseCrashingType.id lacks a db tag"

	_, actual := repo.Read(1)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSQLRepository_ReadGetSqlErr(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := NewSQLRepository[repoTestUser](sqlxDb, -1, "Users", "UserId")

	expected := "unknown dialect"
	_, actual := repo.Read(1)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_ReadPrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM Users WHERE UserId = \\?$").
		WillReturnError(expected)

	_, actual := repo.Read(1)

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_ReadStructScanErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := "missing destination name AnyId in *crudsql.repoTestUser"

	rows := sqlmock.NewRows([]string{"AnyId"}).
		AddRow(1)
	mock.ExpectPrepare("^SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM Users WHERE UserId = \\?$").
		ExpectQuery().
		WithArgs(1).
		WillReturnRows(rows)

	_, actual := repo.Read(1)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_ReadAll(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()

	expected := []repoTestUser{
		{
			ID:        1,
			Name:      "AnyName1",
			Surname:   "AnySurname1",
			Birthdate: time.Now(),
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "AnyName2",
			Surname:   "AnySurname2",
			Birthdate: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"UserId", "Name", "Surname", "Birthdate", "CreatedAt"}).
		AddRow(expected[0].ID, expected[0].Name, expected[0].Surname, expected[0].Birthdate, expected[0].CreatedAt).
		AddRow(expected[1].ID, expected[1].Name, expected[1].Surname, expected[1].Birthdate, expected[1].CreatedAt)
	mock.ExpectPrepare("^SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM Users$").
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

func TestSqlRepository_ReadAllParsePropertiesErr(t *testing.T) {
	repo, mockDB, _, _ := newMock[parseCrashingType]()
	defer mockDB.Close()
	expected := "parseCrashingType.id lacks a db tag"

	_, actual := repo.ReadAll()

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_ReadAllPrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM Users$").
		WillReturnError(expected)

	_, actual := repo.ReadAll()

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_ReadAllStructScanErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := "missing destination name AnyId in *[]crudsql.repoTestUser"

	rows := sqlmock.NewRows([]string{"AnyId"}).
		AddRow(1)
	mock.ExpectPrepare("^SELECT UserId, Name, Surname, Birthdate, CreatedAt FROM Users$").
		ExpectQuery().
		WillReturnRows(rows)

	_, actual := repo.ReadAll()

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_Update(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()

	user := repoTestUser{
		ID:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	mock.ExpectPrepare("^UPDATE Users SET \\(Name = \\?, Surname = \\?, Birthdate = \\?, CreatedAt = \\?\\) WHERE UserId = \\?$").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(1, user)
	if err != nil {
		t.Fatalf("Expected Update to succeed, but got: %s", err)
	}
}

func TestSqlRepository_UpdateParsePropertiesErr(t *testing.T) {
	repo, mockDB, _, _ := newMock[parseCrashingType]()
	defer mockDB.Close()
	expected := "parseCrashingType.id lacks a db tag"

	model := parseCrashingType{}
	actual := repo.Update(1, model)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSQLRepository_UpdateGetSqlErr(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := NewSQLRepository[repoTestUser](sqlxDb, -1, "Users", "UserId")

	expected := "unknown dialect"
	actual := repo.Update(1, repoTestUser{})

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_UpdatePrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^UPDATE Users SET \\(Name = \\?, Surname = \\?, Birthdate = \\?, CreatedAt = \\?\\) WHERE UserId = \\?$").
		WillReturnError(expected)

	actual := repo.Update(1, repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_UpdateExecErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^UPDATE Users SET \\(Name = \\?, Surname = \\?, Birthdate = \\?, CreatedAt = \\?\\) WHERE UserId = \\?$").
		ExpectExec().
		WillReturnError(expected)

	actual := repo.Update(1, repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_UpdateRowsAffectedErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	user := repoTestUser{}
	mock.ExpectPrepare("^UPDATE Users SET \\(Name = \\?, Surname = \\?, Birthdate = \\?, CreatedAt = \\?\\) WHERE UserId = \\?$").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt, 1).
		WillReturnResult(sqlmock.NewErrorResult(expected))

	actual := repo.Update(1, repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_UpdateOtherThanOneRowAffected(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := "2 rows affected by UPDATE statement"
	user := repoTestUser{}
	mock.ExpectPrepare("^UPDATE Users SET \\(Name = \\?, Surname = \\?, Birthdate = \\?, CreatedAt = \\?\\) WHERE UserId = \\?$").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt, 1).
		WillReturnResult(sqlmock.NewResult(1, 2))

	actual := repo.Update(1, repoTestUser{})

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_Delete(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()

	mock.ExpectPrepare("^DELETE FROM Users WHERE UserId = \\?$").
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(1)
	if err != nil {
		t.Fatalf("Expected Delete to succeed, but got: %s", err)
	}
}

func TestSqlRepository_DeletePrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^DELETE FROM Users WHERE UserId = \\?$").
		WillReturnError(expected)

	actual := repo.Delete(1)

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSQLRepository_DeleteGetSqlErr(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := NewSQLRepository[repoTestUser](sqlxDb, -1, "Users", "UserId")

	expected := "unknown dialect"
	actual := repo.Delete(1)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_DeleteExecErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^DELETE FROM Users WHERE UserId = \\?$").
		ExpectExec().
		WillReturnError(expected)

	actual := repo.Delete(1)

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_DeleteRowsAffectedErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("^DELETE FROM Users WHERE UserId = \\?$").
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewErrorResult(expected))

	actual := repo.Delete(1)

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_DeleteOtherThanOneRowAffected(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := "2 rows affected by DELETE statement"
	mock.ExpectPrepare("^DELETE FROM Users WHERE UserId = \\?$").
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 2))

	actual := repo.Delete(1)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestNewSql(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := NewSQLRepository[repoTestUser](sqlxDb, MySQL, "Users", "UserId")

	if repo == nil {
		t.Fatalf("Expected a repo, but got nil instead")
	}
}

func TestNewSqlNilDb(t *testing.T) {
	_, err := NewSQLRepository[repoTestUser](nil, MySQL, "users", "UserId")
	if err == nil {
		t.Fatalf("Expected error on nil db")
	}

	expected := "db cannot be nil"
	if err.Error() != expected {
		t.Fatalf("Expected \"%s\" error but got \"%s\" instead", expected, err.Error())
	}
}

func TestNewSqlEmptyTableName(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	_, err := NewSQLRepository[repoTestUser](sqlxDb, MySQL, "", "UserId")
	if err == nil {
		t.Fatalf("Expected error on empty table name")
	}

	expected := "tableName cannot be empty"
	if err.Error() != expected {
		t.Fatalf("Expected \"%s\" error but got \"%s\" instead", expected, err.Error())
	}
}

func TestNewSqlEmptyIdFieldName(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := NewSQLRepository[repoTestUser](sqlxDb, MySQL, "users", "")
	if repo == nil {
		t.Fatalf("Expected a repo on empty idField")
	}

	if repo.idFieldName != "id" {
		t.Fatalf("Expected idField to be \"id\"")
	}
}
