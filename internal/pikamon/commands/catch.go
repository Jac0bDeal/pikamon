package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/constants"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/models"
	"github.com/bwmarrin/discordgo"
	"github.com/mtslzr/pokeapi-go"
	log "github.com/sirupsen/logrus"
)

var pokemonExpiredMessages = []string{
	"The pokémon heard you coming and ran for the hills!",
	"The pokémon got scared and fled!",
	"The pokémon got away!",
}

func publishExpiredPokemon(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	msg := &discordgo.MessageEmbed{
		Title:       "The pokémon has run away!",
		Description: pokemonExpiredMessages[rand.Intn(len(pokemonExpiredMessages))],
		Color:       constants.MessageColor,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	return err
}

func publichCatchFailure(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	msg := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<@%s> has failed to catch the pokémon!", m.Author.ID),
		Color:       constants.MessageColor,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	return err
}

func publishSuccessfulCatch(s *discordgo.Session, m *discordgo.MessageCreate, pokemon string) (err error) {
	msg := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("Congratulations <@%s>, you caught a %s!", m.Author.ID, strings.Title(pokemon)),
		Color:       constants.MessageColor,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	return err
}

func (h *Handler) catch(s *discordgo.Session, m *discordgo.MessageCreate) {
	// lock a mutex here so that multiple people can't catch the list at the same time.
	h.catchMtx.Lock()
	defer h.catchMtx.Unlock()

	// check if channel id is still in cache. If it is not then there is nothing to catch
	p, exists := h.cache.Channel.Get(m.ChannelID)
	if !exists {
		err := publishExpiredPokemon(s, m)
		if err != nil {
			log.Error(err)
		}
		return
	}
	log.Debug("Pokemon cache still exists, attempting catch...")

	// check if trainer attempting catch is registered, and return register suggestion if not
	trainerID := m.Author.ID
	registered, err := h.isRegistered(trainerID)
	if err != nil {
		log.WithField("trainer", m.Author.ID).Error("Error checking if trainer is registered")
	}
	if !registered {
		publishTrainerNotRegistered(s, m)
		return
	}

	// Create list information object from the cache
	var pokemonId = p.(int)
	pInfo, err := pokeapi.Pokemon(strconv.Itoa(pokemonId))
	if err != nil {
		log.Error(err)
		return
	}

	// Get everything after the "catch" command
	text := strings.TrimSpace(strings.ToLower(m.Content))
	commandText := strings.TrimSpace(text[len(CommandKeyword):])
	commands := strings.Fields(commandText)[1:]
	log.Tracef("Command String: %v\n", commands)

	// Get the list name specified by the person trying to catch it and handle the case where no name was passed
	if len(commands) == 0 {
		err := publichCatchFailure(s, m)
		if err != nil {
			log.Error(err)
		}
		return
	}
	pokemonName := strings.ToLower(commands[0])

	log.WithFields(log.Fields{
		"list":    pokemonName,
		"trainer": trainerID,
	}).Debug("Trying to catch a list!")

	// Perform catch attempt
	expectedPokemonName := pInfo.Name
	caught := strings.EqualFold(pokemonName, expectedPokemonName)

	if caught {
		pokemon := &models.Pokemon{
			PokemonID: pokemonId,
			TrainerID: trainerID,
			Name:      pokemonName,
		}
		if err := h.store.CreatePokemon(pokemon); err != nil {
			log.WithFields(log.Fields{
				"list":    pokemonName,
				"trainer": trainerID,
			}).Error("Failed to create new list in store.")
			return
		}
		err = publishSuccessfulCatch(s, m, expectedPokemonName)
		if err != nil {
			log.Error(err)
		}

		log.WithField("channel", m.ChannelID).Debug("Removing channel from cache")
		h.cache.Channel.Del(m.ChannelID)
	} else {
		// TODO: block same user from trying to immediately re-catch.

		err := publichCatchFailure(s, m)
		if err != nil {
			log.Error(err)
		}
	}
}
