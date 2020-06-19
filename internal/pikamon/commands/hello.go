package commands

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpCommand := "Welcome to the Pikamon bot! The bot currently supports the following commands:\n" +
		"- `p!ka help`\n" +
		"- `p!ka catch <pokemon name> with <pokeball>`\n" +
		"Note: Currently the `with <pokeball>` is not required for the catch."
	msg := discordgo.MessageEmbed{
		Description: helpCommand,
		Color:       0x008080,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Error(err)
	}
}
