package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
)

type Connector struct {
	dbName string
	server string
	port   string
	user   string
	pass   string
}

func NewConnectionManager(dbName, port, server, user, pass string) (Connector, error) {
	return Connector{
		dbName: dbName,
		server: server,
		port:   port,
		user:   user,
		pass:   pass,
	}, nil
}

func (m *Connector) OpenConnect() (*sqlx.DB, error) {
	strConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.user, m.pass, m.server, m.port, m.dbName)


	return m.connect(strConn)
}

func (m *Connector) OpenConnectWithTimeZone(timeZone string) (*sqlx.DB, error) {
	strConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s", m.user, m.pass, m.server, m.port, m.dbName, timeZone)

	return m.connect(strConn)
}

func (m *Connector) connect(strConn string) (*sqlx.DB, error) {
	const (
		_maxIdleCon  = 25
		_maxOpenCon  = 25
		_maxLifeTime = 5
	)

	db, err := sqlx.Connect("mysql", strConn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(_maxIdleCon)
	db.SetMaxOpenConns(_maxOpenCon)
	db.SetConnMaxLifetime(_maxLifeTime)

	return db, nil
}