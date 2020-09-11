package store

import (
	"database/sql"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/config"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	errNoRowsFound = "sql: no rows in result set"
)

// Sqlite is used for interacting with Sqlite as the Store.
type Sqlite struct {
	db       *sql.DB
	location string
}

// NewSqlite builds and returns a pointer to a Sqlite Store implementation from the
// passed config.
func NewSqlite(cfg *config.Config) (*Sqlite, error) {
	return &Sqlite{
		location: cfg.Store.Sqlite.Location,
	}, nil
}

// Close closes the connection to sqlite.
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

// Open opens a connection to sqlite.
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

// CreatePokemon creates a Pokemon assigned to a Trainer.
func (s *Sqlite) CreatePokemon(pokemon *models.Pokemon) error {
	stmt, err := s.db.Prepare(
		"INSERT INTO pokemon (trainer, pokemon_id, name) VALUES (?, ?, ?)",
	)
	if err != nil {
		return errors.Wrap(err, "could not prepare pokemon insert statement")
	}
	if _, err := stmt.Exec(
		pokemon.TrainerID,
		pokemon.PokemonID,
		pokemon.Name,
	); err != nil {
		return errors.Wrap(err, "failed to execute pokemon insert statement against store")
	}
	return nil
}

// CreateTrainer creates a new Trainer.
func (s *Sqlite) CreateTrainer(trainer *models.Trainer) error {
	stmt, err := s.db.Prepare(
		"INSERT INTO trainers (id) VALUES (?)",
	)
	if err != nil {
		return errors.Wrap(err, "could not prepare trainer insert statement")
	}
	if _, err := stmt.Exec(trainer.ID); err != nil {
		return errors.Wrap(err, "failed to execute trainer insert statement against store")
	}
	return nil
}

// GetAllPokemon returns all of the pokemon currently assigned to a trainer.
func (s *Sqlite) GetAllPokemon(trainerID string) ([]*models.Pokemon, error) {
	stmt, err := s.db.Prepare(
		"SELECT id, pokemon_id, name FROM pokemon WHERE trainer = ?",
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare trainer select statement")
	}

	rows, err := stmt.Query(trainerID)
	if err != nil {
		if err.Error() == errNoRowsFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to execute trainer select statement against store")
	}

	var pokemon []*models.Pokemon
	for rows.Next() {
		var p models.Pokemon
		p.TrainerID = trainerID
		err = rows.Scan(&p.ID, &p.PokemonID, &p.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to process row returned from trainer select query")
		}
		pokemon = append(pokemon, &p)
	}
	return pokemon, nil
}

// GetTrainer returns the information associated with a trainer.
func (s *Sqlite) GetTrainer(trainerID string) (*models.Trainer, error) {
	stmt, err := s.db.Prepare(
		"SELECT id FROM trainers WHERE id = ?",
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare trainer select statement")
	}
	var trainer models.Trainer
	err = stmt.QueryRow(trainerID).Scan(&trainer.ID)
	if err != nil {
		if err.Error() == errNoRowsFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to execute trainer select statement against store")
	}
	return &trainer, nil
}
