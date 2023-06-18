package postgres

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"

	"Cryptogo/config"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPsqlDB(cfg *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Postgres.PostgreSQLHost,
		cfg.Postgres.PostgreSQLPort,
		cfg.Postgres.PostgreSQLUser,
		cfg.Postgres.PostgreSQLDBName,
		cfg.Postgres.PostgreSQLPassword,
	)

	db, err := sqlx.Connect(cfg.Postgres.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
