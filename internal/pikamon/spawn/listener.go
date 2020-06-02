package spawn

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
)

type spawner interface {
	spawn(*discordgo.Session, *discordgo.MessageCreate) bool
}

// Handler listens to non-Pikamon messages in a channel, calls performs the
// spawn operations.
type Handler struct {
	channelCache   *ristretto.Cache
	debounceWindow time.Duration
	spawners       []spawner
}

// NewHandler constructs and returns a new Handler that spawns things in channels.
func NewHandler(pokemonSpawnChance float64, debounceWindow time.Duration) (*Handler, error) {
	pokemonSpawner, err := newPokemonSpawner(pokemonSpawnChance, debounceWindow)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build pokemon spawner")
	}
	return &Handler{
		spawners: []spawner{
			pokemonSpawner,
		},
	}, nil
}

// Handle is the handler function registered on the discord bot that
// processes incoming messages and calls into each spawner.
func (h *Handler) Handle(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == sess.State.User.ID {
		return
	}

	for _, s := range h.spawners {
		s.spawn(sess, m)
	}
}
