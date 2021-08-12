package commands

import (
	"fmt"
	"strconv"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/constants"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func publishTrainerAlreadyRegistered(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.WithFields(log.Fields{"id": m.Author.ID}).Debugf("Trainer is already registered")
	msg := &discordgo.MessageEmbed{
		Title:       "You've already registered, silly!",
		Description: fmt.Sprintf("For future reference, your trainer ID is %s.", m.Author.ID),
		Color:       constants.MessageColor,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, msg); err != nil {
		log.Error(err)
	}
}

func publishTrainerWelcome(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.WithFields(log.Fields{"id": m.Author.ID}).Debugf("Welcoming new trainer!")
	msg := &discordgo.MessageEmbed{
		Title:       "Welcome to Pikamon!",
		Description: fmt.Sprintf("Thanks for registering, <@%s>! Your trainer ID is %s.", m.Author.ID, m.Author.ID),
		Color:       constants.MessageColor,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, msg); err != nil {
		log.Error(err)
	}
}

func (h *Handler) register(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.WithFields(log.Fields{"id": m.Author.ID}).Debug("Received init trainer command")
	trainerID, err := strconv.Atoi(m.Author.ID)
	if err != nil {
		log.Errorf("Received unparsable integer %s: %v", m.Author.ID, err)
	}

	log.WithFields(log.Fields{"id": trainerID}).Debugf("Checking if trainer is already registered...")
	registered, trainerID, err := h.isRegistered(m.Author.ID)
	if err != nil {
		log.WithField("trainer", m.Author.ID).Error("Error checking if trainer is registered")
	}
	if registered {
		publishTrainerAlreadyRegistered(s, m)
		return
	}

	log.WithFields(log.Fields{"id": trainerID}).Info("Trainer not found, registering now...")
	newTrainer := &models.Trainer{
		ID: trainerID,
	}
	if err = h.store.CreateTrainer(newTrainer); err != nil {
		log.WithFields(log.Fields{"id": trainerID}).Error("Failed to create trainer in store: %v", err)
	}

	publishTrainerWelcome(s, m)
	log.WithFields(log.Fields{"id": trainerID}).Debugf("Trainer registered successfully")
}
