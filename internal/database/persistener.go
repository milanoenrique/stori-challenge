package database

import (
	"github.com/jmoiron/sqlx"
)

const (
	DatabaseTimeLayout = "2006-01-02 15:04:05"
)

type Persistence struct {
	DBConector DBConector
}

type DBConector interface {
	OpenConnect() (*sqlx.DB, error)
	OpenConnectWithTimeZone(timeZone string) (*sqlx.DB, error)
}

func NewPersistence(theDBConector DBConector) Persistence {
	return Persistence{
		DBConector: theDBConector,
	}
}