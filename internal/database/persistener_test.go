package database_test

import (
	"payment-process/internal/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/tj/assert"
)

type stubDBManager struct {
	db *sqlx.DB
}

func (d *stubDBManager) OpenConnectWithTimeZone(timeZone string) (*sqlx.DB, error) {
	return d.db, nil
}


func (d *stubDBManager) OpenConnect() (*sqlx.DB, error) {
	return d.db, nil
}

func TestNewPersistenct(t *testing.T) {
	mockDB, _, _ := sqlmock.New()

	stub := stubDBManager{
		db: sqlx.NewDb(mockDB, "sqlmock"),
	}
	p := database.NewPersistence(&stub)
	assert.NotNil(t, p)
}
