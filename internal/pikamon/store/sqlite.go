package store

import (
	"database/sql"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/config"
	log "github.com/sirupsen/logrus"

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
	log.Info("Closing sqlite db...")
	if s.db == nil {
		log.Info("DB already closed, nothing to do.")
		return nil
	}
	if err := s.db.Close(); err != nil {
		return err
	}
	log.Info("Closed sqlite db.")
	return nil
}

func (s *Sqlite) Open() error {
	log.Infof("Opening sqlite db at '%s'...", s.location)
	db, err := sql.Open("sqlite3", s.location)
	if err != nil {
		return errors.Wrapf(err, "failed to open sqlite db at '%s'", s.location)
	}
	s.db = db
	log.Info("Opened sqlite db.")
	return nil
}
