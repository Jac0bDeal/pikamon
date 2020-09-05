package store

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/config"
	"github.com/pkg/errors"
)

const (
	TypeSqlite = "sqlite"
)

// Store represents the various operations for interacting with a Pikamon Store.
type Store interface {
	Open() error
	Close() error
}

// New builds and returns an implementation of Store based on the passed config.
func New(cfg *config.Config) (Store, error) {
	switch cfg.Store.Type {
	case TypeSqlite:
		return newSqlite(cfg)
	default:
		return nil, errors.New("unrecognized store type")
	}
}
