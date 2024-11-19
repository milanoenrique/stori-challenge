package repositories_test

import (

	"errors"
	"fmt"
	"payment-process/internal/database"
	"payment-process/internal/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/tj/assert"
)

type stubDBManager struct {
	db *sqlx.DB
}

// OpenConnect implements database.DBConector.
func (s *stubDBManager) OpenConnect() (*sqlx.DB, error) {
	return s.db, nil
}

// OpenConnectWithTimeZone implements database.DBConector.
func (s *stubDBManager) OpenConnectWithTimeZone(timeZone string) (*sqlx.DB, error) {
	panic("unimplemented")
}

func TestInsertTransaction_success(t *testing.T) {

	// Arrange
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.FailNow()
	}

	mock.ExpectExec("INSERT INTO transactions").WithArgs(0, "123456", "8/11", 100.00).WillReturnResult(sqlmock.NewResult(1, 1))

	stub := stubDBManager{
		db: sqlx.NewDb(mockDB, "sqlmock"),
	}

	p := database.NewPersistence(&stub)
	repo := repositories.NewTXRepository(p)
	tx := repositories.Transaction{}
	tx.AccountId = "123456"
	tx.Amount = 100
	tx.Date = "8/11"
	tx.Id = 0
	err = repo.InsertTransaction(&tx)

	assert.Nil(t, err)
}

func TestInsertTransaction_fail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.FailNow()
	}

	mock.ExpectExec("INSERT INTO transactions").WithArgs(0, "123456", "8/11", 100.00).WillReturnError(errors.New(""))

	stub := stubDBManager{
		db: sqlx.NewDb(mockDB, "sqlmock"),
	}

	p := database.NewPersistence(&stub)
	repo := repositories.NewTXRepository(p)
	tx := repositories.Transaction{}
	tx.AccountId = "123456"
	tx.Amount = 100
	tx.Date = "8/11"
	tx.Id = 0
	err = repo.InsertTransaction(&tx)

	assert.NotNil(t, err)
}

func TestGetTransactions_success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.FailNow()
	}

	rows := sqlmock.NewRows([]string{"id", "account_id", "date", "amount"}).AddRow(0, "123456","07/11", "100.00")

	mock.ExpectQuery(`^SELECT (.+) account_id = ?`).WithArgs("123456").WillReturnRows(rows)

	stub := stubDBManager{
		db: sqlx.NewDb(mockDB, "sqlmock"),
	}

	p := database.NewPersistence(&stub)
	repo := repositories.NewTXRepository(p)
	result, err := repo.GetTransactions("123456")
	if err != nil {
		fmt.Println(err)
	}
	assert.NotEmpty(t, result)
	assert.Equal(t, "123456", result[0].AccountId)
	assert.Equal(t, 100.00, result[0].Amount)
	assert.Equal(t, "07/11", result[0].Date)
}

func TestGetTRansactions_fail(t *testing.T){
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.FailNow()
	}

	mock.ExpectQuery(`^SELECT (.+) account_id = ?`).WithArgs("123456").WillReturnError(errors.New(""))

	stub := stubDBManager{
		db: sqlx.NewDb(mockDB, "sqlmock"),
	}

	p := database.NewPersistence(&stub)
	repo := repositories.NewTXRepository(p)
	result, err := repo.GetTransactions("123456")
	assert.Nil(t, result)
	assert.NotNil(t, err)
}