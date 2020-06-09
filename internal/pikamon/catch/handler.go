package catch

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type catcher interface {
	catch(*discordgo.Session, *discordgo.MessageCreate) bool
}

// Handler listens to Pikamon catch messages in a channel and attempts
// the catch operation with the specified pokeball.
type Handler struct {
	catchers []catcher
}

// TODO - first just build the catch. have it identify a catch message and print back to the channel "You caught it"

// CatchHandler constructs and returns a new Handler that catches pokemon in the channels.
// TODO - figure out if this should take the channel cache from the spawner
func CatchHandler(s *discordgo.Session, m *discordgo.MessageCreate) (*Handler, error) {
	return &Handler{}, nil
}

// CatchHandle is the handler function registered on the discord bot that
// processes incoming messages and checks if we are catching a spawned pokemon
func (h *Handler) CatchHandle(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == sess.State.User.ID {
		return
	}

	fmt.Printf("Author ID: %s", m.Author.ID)
}

/*
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
*/
