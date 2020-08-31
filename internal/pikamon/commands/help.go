package commands

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/constants"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpCommand := "Welcome to the Pikamon bot! The bot currently supports the following commands:\n" +
		"- `p!ka help`\n" +
		"- `p!ka catch <pokemon name>`"
	msg := discordgo.MessageEmbed{
		Description: helpCommand,
		Color:       constants.MessageColor,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Error(err)
	}
}
