package store

import (
	"database/sql"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/config"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Sqlite struct {
	db       *sql.DB
	location string
}

func newSqlite(cfg *config.Config) (*Sqlite, error) {
	return &Sqlite{
		location: cfg.Store.Sqlite.Location,
	}, nil
}

func (s *Sqlite) Close() error {
	// db already closed, nothing to do
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Sqlite) Open() error {
	db, err := sql.Open("sqlite3", s.location)
	if err != nil {
		return errors.Wrapf(err, "failed to open sqlite db at '%s'", s.location)
	}
	s.db = db
	return nil
}
