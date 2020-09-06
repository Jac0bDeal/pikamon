package store

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/config"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/models"

	"github.com/pkg/errors"
)

const (
	TypeSqlite = "sqlite"
)

// Store represents the various operations for interacting with a Pikamon Store.
type Store interface {
	Open() error
	Close() error

	CreatePokemon(pokemon *models.Pokemon) error
	CreateTrainer(trainer *models.Trainer) error
	GetAllPokemon(trainer int) ([]*models.Pokemon, error)
	GetTrainer(trainer int) (*models.Trainer, error)
}

// New builds and returns an implementation of Store based on the passed config.
func New(cfg *config.Config) (Store, error) {
	switch cfg.Store.Type {
	case TypeSqlite:
		return NewSqlite(cfg)
	default:
		return nil, errors.New("unrecognized store type")
	}
}
