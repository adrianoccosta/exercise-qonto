package config

import (
	"database/sql"
	"github.com/adrianoccosta/exercise-qonto/log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Conn Contains information for current db connection.
type Conn struct {
	Conn *sql.DB
}

// DBConnection DB connections mandatory fields
type DBConnection struct {
	Path              string
	MaxIdleConns      int
	MaxOpenConns      int
	ConnMaxTTLMinutes int
}

// InitDBConnection function init the database connection.
func InitDBConnection(conn DBConnection, logger log.Logger) Conn {
	db, err := sql.Open("sqlite3", conn.Path)

	if err != nil {
		if err != nil {
			logger.Fatal("cannot be instanced without an db instance")
		}
	}

	db.SetMaxIdleConns(conn.MaxIdleConns)
	db.SetMaxOpenConns(conn.MaxOpenConns)
	db.SetConnMaxLifetime(time.Minute * time.Duration(conn.ConnMaxTTLMinutes))

	return Conn{Conn: db}
}
