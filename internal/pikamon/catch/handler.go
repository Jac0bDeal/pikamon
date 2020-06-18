package catch

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/util"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	CommandKeyword = "p!ka"
)

type catcher interface {
	catch(*discordgo.Session, *discordgo.MessageCreate) bool
}

// Handler listens to non-Pikamon messages in a channel and  performs the
// spawn operations.
type Handler struct {
	catchers []catcher
}

// Handler listens to Pikamon catch messages in a channel and attempts
// the catch operation with the specified pokeball.
type pokemonCatcher struct {
	cache *util.BotCache
}

// NewHandler constructs and returns a new Handler that spawns things in channels.
func NewCatchHandler(botCache *util.BotCache) (*Handler, error) {
	catchPokemon, err := newPokemonCatcher(botCache)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build pokemon spawner")
	}

	return &Handler{
		catchers: []catcher{
			catchPokemon,
		},
	}, nil
}

func newPokemonCatcher(botCache *util.BotCache) (*pokemonCatcher, error) {
	return &pokemonCatcher{
		cache: botCache,
	}, nil
}

// Handle is the handler function registered on the discord bot that
// processes incoming messages and calls into each spawner.
func (h *Handler) Handle(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == sess.State.User.ID {
		return
	}

	for _, s := range h.catchers {
		s.catch(sess, m)
	}
}

func (c *pokemonCatcher) catch(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	text := strings.TrimSpace(strings.ToLower(m.Content))
	commandText := strings.TrimSpace(text[len(CommandKeyword):])

	// Get everything after the "catch" command
	commands := strings.Fields(commandText)[1:]
	log.Info("Command String: %v\n", commands)

	pokemon := strings.ToLower(commands[0])

	// Check to see if they specify a pokeball type
	var pokeball string
	if len(commands) > 1 && strings.ToLower(commands[1]) == "with" {
		pokeball = strings.ToLower(commands[2])
	}

	log.WithFields(log.Fields{
		"pokemon":  pokemon,
		"pokeball": pokeball,
	}).Info("Trying to catch a pokemon!")
}
