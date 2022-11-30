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
	repo, _ := New[T](sqlxDb, MySQL, "Users", "UserId")
	repo.templates = sqlTemplatesMock{}
	return repo, mockDB, mock, err
}

func TestSqlRepository_Create(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetInsertMock: func(fields []string) (string, error) {
			return "AnyInsert", nil
		},
	}

	user := repoTestUser{
		ID:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	mock.ExpectPrepare("AnyInsert").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Expected Create to succeed, but got: %s", err)
	}
}

func TestSQLRepository_CreateGetSqlErr(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := New[repoTestUser](sqlxDb, -1, "Users", "UserId")
	expected := fmt.Errorf("AnyError")
	repo.templates = sqlTemplatesMock{
		GetInsertMock: func(fields []string) (string, error) {
			return "", expected
		},
	}

	actual := repo.Create(repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_CreatePrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetInsertMock: func(fields []string) (string, error) {
			return "AnyInsert", nil
		},
	}
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("AnyInsert").
		WillReturnError(expected)

	actual := repo.Create(repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_CreateExecErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetInsertMock: func(fields []string) (string, error) {
			return "AnyInsert", nil
		},
	}
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("AnyInsert").
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
	repo.templates = sqlTemplatesMock{
		GetInsertMock: func(fields []string) (string, error) {
			return "AnyInsert", nil
		},
	}
	expected := fmt.Errorf("any error")
	user := repoTestUser{}
	mock.ExpectPrepare("AnyInsert").
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
	repo.templates = sqlTemplatesMock{
		GetInsertMock: func(fields []string) (string, error) {
			return "AnyInsert", nil
		},
	}
	expected := "2 rows affected by INSERT INTO statement"
	user := repoTestUser{}
	mock.ExpectPrepare("AnyInsert").
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
	repo.templates = sqlTemplatesMock{
		GetSelectMock: func() string {
			return "AnySelect"
		},
	}

	expected := repoTestUser{
		ID:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"UserId", "Name", "Surname", "Birthdate", "CreatedAt"}).
		AddRow(expected.ID, expected.Name, expected.Surname, expected.Birthdate, expected.CreatedAt)
	mock.ExpectPrepare("AnySelect").
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

func TestSqlRepository_ReadPrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetSelectMock: func() string {
			return "AnySelect"
		},
	}
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("AnySelect").
		WillReturnError(expected)

	_, actual := repo.Read(1)

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_ReadStructScanErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetSelectMock: func() string {
			return "AnySelect"
		},
	}
	expected := "missing destination name AnyId in *crudsql.repoTestUser"

	rows := sqlmock.NewRows([]string{"AnyId"}).
		AddRow(1)
	mock.ExpectPrepare("AnySelect").
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
	repo.templates = sqlTemplatesMock{
		GetSelectAllMock: func() string {
			return "AnySelectAll"
		},
	}

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
	mock.ExpectPrepare("AnySelectAll").
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

func TestSqlRepository_ReadAllPrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetSelectAllMock: func() string {
			return "AnySelectAll"
		},
	}
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("AnySelectAll").
		WillReturnError(expected)

	_, actual := repo.ReadAll()

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_ReadAllStructScanErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetSelectAllMock: func() string {
			return "AnySelectAll"
		},
	}
	expected := "missing destination name AnyId in *[]crudsql.repoTestUser"

	rows := sqlmock.NewRows([]string{"AnyId"}).
		AddRow(1)
	mock.ExpectPrepare("AnySelectAll").
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
	repo.templates = sqlTemplatesMock{
		GetUpdateMock: func(fields []string) (string, error) {
			return "AnyGetUpdate", nil
		},
	}

	user := repoTestUser{
		ID:        1,
		Name:      "AnyName",
		Surname:   "AnySurname",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
	}

	mock.ExpectPrepare("AnyGetUpdate").
		ExpectExec().
		WithArgs(user.Name, user.Surname, user.Birthdate, user.CreatedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(1, user)
	if err != nil {
		t.Fatalf("Expected Update to succeed, but got: %s", err)
	}
}

func TestSQLRepository_UpdateGetSqlErr(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	expected := fmt.Errorf("AnyError")
	repo, _ := New[repoTestUser](sqlxDb, -1, "Users", "UserId")
	repo.templates = sqlTemplatesMock{
		GetUpdateMock: func(fields []string) (string, error) {
			return "", expected
		},
	}

	actual := repo.Update(1, repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_UpdatePrepareErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	expected := fmt.Errorf("any error")
	repo.templates = sqlTemplatesMock{
		GetUpdateMock: func(fields []string) (string, error) {
			return "AnyUpdate", nil
		},
	}
	mock.ExpectPrepare("AnyUpdate").
		WillReturnError(expected)

	actual := repo.Update(1, repoTestUser{})

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_UpdateExecErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetUpdateMock: func(fields []string) (string, error) {
			return "AnyUpdate", nil
		},
	}
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("AnyUpdate").
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
	repo.templates = sqlTemplatesMock{
		GetUpdateMock: func(fields []string) (string, error) {
			return "AnyUpdate", nil
		},
	}
	expected := fmt.Errorf("any error")
	user := repoTestUser{}
	mock.ExpectPrepare("AnyUpdate").
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
	repo.templates = sqlTemplatesMock{
		GetUpdateMock: func(fields []string) (string, error) {
			return "AnyUpdate", nil
		},
	}
	expected := "2 rows affected by UPDATE statement"
	user := repoTestUser{}
	mock.ExpectPrepare("AnyUpdate").
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
	repo.templates = sqlTemplatesMock{
		GetDeleteMock: func() string {
			return "AnyDelete"
		},
	}

	mock.ExpectPrepare("AnyDelete").
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
	repo.templates = sqlTemplatesMock{
		GetDeleteMock: func() string {
			return "AnyDelete"
		},
	}
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("AnyDelete").
		WillReturnError(expected)

	actual := repo.Delete(1)

	if actual != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestSqlRepository_DeleteExecErr(t *testing.T) {
	repo, mockDB, mock, _ := newMock[repoTestUser]()
	defer mockDB.Close()
	repo.templates = sqlTemplatesMock{
		GetDeleteMock: func() string {
			return "AnyDelete"
		},
	}
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("AnyDelete").
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
	repo.templates = sqlTemplatesMock{
		GetDeleteMock: func() string {
			return "AnyDelete"
		},
	}
	expected := fmt.Errorf("any error")
	mock.ExpectPrepare("AnyDelete").
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
	repo.templates = sqlTemplatesMock{
		GetDeleteMock: func() string {
			return "AnyDelete"
		},
	}
	expected := "2 rows affected by DELETE statement"
	mock.ExpectPrepare("AnyDelete").
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 2))

	actual := repo.Delete(1)

	if actual.Error() != expected {
		t.Fatalf("Expected \"%s\" but got \"%s\" instead", expected, actual)
	}
}

func TestNew(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := New[repoTestUser](sqlxDb, MySQL, "Users", "UserId")

	if repo == nil {
		t.Fatalf("Expected a repo, but got nil instead")
	}
}

func TestNewNilDb(t *testing.T) {
	_, err := New[repoTestUser](nil, MySQL, "users", "UserId")
	if err == nil {
		t.Fatalf("Expected error on nil db")
	}

	expected := "db cannot be nil"
	if err.Error() != expected {
		t.Fatalf("Expected \"%s\" error but got \"%s\" instead", expected, err.Error())
	}
}

func TestNewEmptyTableName(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	_, err := New[repoTestUser](sqlxDb, MySQL, "", "UserId")
	if err == nil {
		t.Fatalf("Expected error on empty table name")
	}

	expected := "tableName cannot be empty"
	if err.Error() != expected {
		t.Fatalf("Expected \"%s\" error but got \"%s\" instead", expected, err.Error())
	}
}

func TestNewEmptyIdFieldName(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDb := sqlx.NewDb(mockDB, "sqlmock")
	repo, _ := New[repoTestUser](sqlxDb, MySQL, "users", "")
	if repo == nil {
		t.Fatalf("Expected a repo on empty idField")
	}

	if repo.idFieldName != "id" {
		t.Fatalf("Expected idField to be \"id\"")
	}
}
