package commands

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/util"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type handler func(*discordgo.Session, *discordgo.MessageCreate)

var rootCommands = map[string]handler{
	"help":  help,
	"catch": catch,
}

// Handle processes all the commands that match the bot command keyword; chaining
// handlers until there isn't a recognized command, an error occurs, or everything
// succeeds.
func Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Tracef("Received message create: %+v", m.Message)

	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	text := strings.TrimSpace(strings.ToLower(m.Content))

	// ignore all messages not prefixed with bot command keyword
	if !strings.HasPrefix(text, util.CommandKeyword) {
		return
	}
	commandText := strings.TrimSpace(text[len(util.CommandKeyword):])
	if commandText == "" {
		return
	}

	// call the appropriate handler based on the root command
	commands := strings.Fields(commandText)

	command := "help" // Default to help if not enough commands are passed
	if len(commands) > 0 {
		command = commands[0]
	}

	next, exists := rootCommands[command]
	if !exists {
		log.Debugf("Received unrecognized command: %s", commandText)
		next = rootCommands["help"] // Default to help for invalid command
	}

	log.Debugf("Received command: %s", commandText)
	next(s, m)
}
