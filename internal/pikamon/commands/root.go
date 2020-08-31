package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	log "github.com/sirupsen/logrus"
)

// Pikamon bot keyword
const (
	CommandKeyword = "p!ka"

	catchCommand = "catch"
	helpCommand  = "help"
)

type handler func(*discordgo.Session, *discordgo.MessageCreate)

type Handler struct {
	channelCache *ristretto.Cache
}

func NewHandler(channelCache *ristretto.Cache) *Handler {
	return &Handler{
		channelCache: channelCache,
	}
}

// Handle processes all the commands that match the bot command keyword; chaining
// handlers until there isn't a recognized command, an error occurs, or everything
// succeeds.
func (h *Handler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Tracef("Received message create: %+v", m.Message)

	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	text := strings.TrimSpace(strings.ToLower(m.Content))

	// ignore all messages not prefixed with bot command keyword
	if !strings.HasPrefix(text, CommandKeyword) {
		return
	}
	commandText := strings.TrimSpace(text[len(CommandKeyword):])
	if commandText == "" {
		return
	}

	// call the appropriate handler based on the root command
	commands := strings.Fields(commandText)

	command := helpCommand // Default to help if not enough commands are passed
	if len(commands) > 0 {
		command = commands[0]
	}

	var handle handler
	switch command {
	case catchCommand:
		handle = h.catch
	case helpCommand:
		handle = h.help
	default:
		log.Debugf("Received unrecognized command: %s", commandText)
		handle = h.help
	}

	log.Debugf("Received command: %s", commandText)
	handle(s, m)
}
