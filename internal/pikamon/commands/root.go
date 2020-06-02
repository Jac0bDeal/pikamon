package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	commandKeyword = "p!ka"
)

type handler func(*discordgo.Session, *discordgo.MessageCreate)

var rootCommands = map[string]handler{
	"help": help,
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
	if !strings.HasPrefix(text, commandKeyword) {
		return
	}
	commandText := strings.TrimSpace(text[len(commandKeyword):])
	if commandText == "" {
		return
	}

	// call the appropriate handler based on the root command
	command := strings.SplitN(commandText, " ", 1)[0]
	next, exists := rootCommands[command]
	if !exists {
		log.Debugf("Received unrecognized command: %s", commandText)
		return
	}
	log.Debugf("Received command: %s", commandText)
	next(s, m)
}
