package commands

import (
	"strings"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/constants"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) list(s *discordgo.Session, m *discordgo.MessageCreate) {
	trainerID := m.Author.ID
	// check if trainer is registered, and return register suggestion if not
	registered, err := h.isRegistered(trainerID)
	if err != nil {
		log.WithField("trainer", trainerID).Error("Error checking if trainer is registered")
	}
	if !registered {
		publishTrainerNotRegistered(s, m)
		return
	}

	pokemon, err := h.store.GetAllPokemon(trainerID)
	if err != nil {
		log.WithField("trainer", trainerID).Errorf("Failed to get all list for trainer: %v", err)
	}

	publishPokemonInfo(pokemon, s, m)
}

func publishPokemonInfo(pokemon []*models.Pokemon, s *discordgo.Session, m *discordgo.MessageCreate) {
	pokemonInfo := make([]string, len(pokemon))
	for idx, p := range pokemon {
		pokemonInfo[idx] = p.ListingInfo()
	}
	msg := &discordgo.MessageEmbed{
		Title:       "Your pokémon:",
		Description: strings.Join(pokemonInfo, "\n"),
		Color:       constants.MessageColor,
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		log.Error(err)
	}
}
