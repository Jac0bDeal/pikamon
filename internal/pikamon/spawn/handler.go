package spawn

import (
	"strings"
	"time"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
)

type spawner interface {
	spawn(*discordgo.Session, *discordgo.MessageCreate) bool
}

// Handler listens to non-Pikamon messages in a channel and  performs the
// spawn operations.
type Handler struct {
	spawners []spawner
}

// NewHandler constructs and returns a new Handler that spawns things in channels.
func NewHandler(
	channelCache *ristretto.Cache,
	pokemonSpawnChance float64,
	maximumSpawnDuration time.Duration,
	maxPokemonID int,
) (*Handler, error) {
	pokemonSpawner, err := newPokemonSpawner(channelCache, pokemonSpawnChance, maximumSpawnDuration, maxPokemonID)
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

	// ignore all messages prefixed with bot command keyword
	text := strings.TrimSpace(strings.ToLower(m.Content))
	if strings.HasPrefix(text, commands.CommandKeyword) {
		return
	}

	for _, s := range h.spawners {
		s.spawn(sess, m)
	}
}
