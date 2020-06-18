package commands

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/util"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
	//"github.com/Jac0bDeal/pikamon/internal/pikamon/items"
)

func catch(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Infof("Bot sample value: %s", util.BotMetadata.Sample)

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
