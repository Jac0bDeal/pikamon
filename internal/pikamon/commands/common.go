package commands

import (
	"fmt"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/constants"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) isRegistered(trainerID string) (bool, error) {
	registeredTrainer, err := h.store.GetTrainer(trainerID)
	if err != nil {
		return false, errors.Wrap(err, "unable to fetch trainer from store")
	}
	if registeredTrainer == nil {
		return false, nil
	}
	return true, nil
}

func publishTrainerNotRegistered(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.WithFields(log.Fields{"id": m.Author.ID}).Debugf("Trainer is not registered.")
	msg := &discordgo.MessageEmbed{
		Title:       "You aren't a registered trainer!",
		Description: fmt.Sprintf("Please register using %s %s.", CommandKeyword, registerCommand),
		Color:       constants.MessageColor,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, msg); err != nil {
		log.Error(err)
	}
}
