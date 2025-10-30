package repositories

import (
	"fmt"

	"github.com/warmhouse/warmhouse_devices/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PgDriver struct {
	db *sqlx.DB
}

func NewPgDriver(conf *config.Secrets) (*PgDriver, error) {
	db, err := sqlx.Connect("postgres", conf.Pg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PgDriver{db: db}, nil
}

func (d *PgDriver) DB() *sqlx.DB {
	return d.db
}
