package commands

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/constants"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpCommand := "Welcome to the Pikamon bot! The bot currently supports the following commands:\n" +
		"- `p!ka register` (registers you as a trainer!)\n" +
		"- `p!ka list` (lists your pokémon)\n" +
		"- `p!ka help (shows this help text)`\n" +
		"- `p!ka catch <pokémon name>` (attempt to catch a pokémon by spelling its name correctly)"
	msg := discordgo.MessageEmbed{
		Description: helpCommand,
		Color:       constants.MessageColor,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Error(err)
	}
}
