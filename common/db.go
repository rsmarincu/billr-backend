package common

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	EnvKeyDatabaseConnection = "DATABASE_URL"
)

type DatabaseService interface {
	GetDb() *sql.DB
}

func NewDatabaseService(config Config) DatabaseService {
	connStr := config.GetEnvVariable(EnvKeyDatabaseConnection, "")
	if connStr == "" {
		panic("please provide database connection")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &databaseService{
		db: db,
	}
}

type databaseService struct {
	db *sql.DB
}

func (d *databaseService) GetDb() *sql.DB {
	return d.db
}
