package repositories_test

import (
	"payment-process/internal/database"
	"payment-process/internal/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/tj/assert"
)

func TestFindUserByAccountId(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.FailNow()
	}

	rows := sqlmock.NewRows([]string{"name", "last_name", "email", "account_id"}).AddRow("test", "test", "test@e,ail.com", "123456")

	mock.ExpectQuery(`^SELECT (.+) account_id = ?`).WithArgs("123456").WillReturnRows(rows)

	stub := stubDBManager{
		db: sqlx.NewDb(mockDB, "sqlmock"),
	}

	p := database.NewPersistence(&stub)
	repo := repositories.NewURepository(p)
	result, err:= repo.FindUserByAccountId("123456")

	assert.NoError(t, err)
	assert.NotNil(t, result)

}