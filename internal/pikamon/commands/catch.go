package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/constants"
	"github.com/bwmarrin/discordgo"
	"github.com/mtslzr/pokeapi-go"
	log "github.com/sirupsen/logrus"
)

var pokemonExpiredMessages = []string{
	"The pokemon heard you coming and ran for the hills!",
	"The pokemon got scared and fled!",
	"The pokemon got away!",
}

var catchFailureMessages = []string{
	"Oof, that is the wrong pokemon!",
	"Darn, the pokemon broke free!",
}

func publishExpiredPokemon(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	expireMessageIndex := rand.Intn(len(pokemonExpiredMessages))
	expireMessage := pokemonExpiredMessages[expireMessageIndex]
	msg := discordgo.MessageEmbed{
		Title:       "The Pokemon has run away!",
		Description: expireMessage,
		Color:       constants.MessageColor,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msg)
	return err
}

func publichCatchFailure(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	catchFailureMessageIndex := rand.Intn(len(catchFailureMessages))
	expireMessage := fmt.Sprintf("<@%s> has failed to catch the pokemon! %s", m.Author.ID, catchFailureMessages[catchFailureMessageIndex])
	msg := discordgo.MessageEmbed{
		Description: expireMessage,
		Color:       constants.MessageColor,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msg)
	return err
}

func publishSuccessfulCatch(s *discordgo.Session, m *discordgo.MessageCreate, pokemon string) (err error) {
	// TODO - save to database
	catchMessage := fmt.Sprintf("Congratulations <@%s>, you caught a %s!", m.Author.ID, strings.Title(pokemon))

	msg := discordgo.MessageEmbed{
		Description: catchMessage,
		Color:       constants.MessageColor,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msg)
	return err
}

func (h *Handler) catch(s *discordgo.Session, m *discordgo.MessageCreate) {
	h.catchMtx.Lock()
	defer h.catchMtx.Unlock()
	// check if channel id is still in cache. If it is not then there is nothing to catch
	p, exists := h.channelCache.Get(m.ChannelID)
	if !exists {
		err := publishExpiredPokemon(s, m)
		if err != nil {
			log.Error(err)
		}
		return
	}
	log.Debug("Pokemon cache still exists, attempting catch...")

	// Create pokemon information object from the cache
	var pokemonId = p.(int)
	pInfo, err := pokeapi.Pokemon(strconv.Itoa(pokemonId))
	if err != nil {
		log.Error(err)
	}

	// Get everything after the "catch" command
	text := strings.TrimSpace(strings.ToLower(m.Content))
	commandText := strings.TrimSpace(text[len(CommandKeyword):])
	commands := strings.Fields(commandText)[1:]
	log.Infof("Command String: %v\n", commands)

	// Get the pokemon name specified by the person trying to catch it and handle the case where no name was passed
	if len(commands) == 0 {
		err := publichCatchFailure(s, m)
		if err != nil {
			log.Error(err)
		}
		return
	}
	pokemonName := strings.ToLower(commands[0])

	// Check to see if they specify a pokeball type
	var pokeball string
	if len(commands) > 1 && strings.ToLower(commands[1]) == "with" {
		pokeball = strings.ToLower(commands[2])
	}

	log.WithFields(log.Fields{
		"pokemon":  pokemonName,
		"pokeball": pokeball,
	}).Debug("Trying to catch a pokemon!")

	// Perform catch attempt
	expectedPokemonName := pInfo.Name
	if strings.EqualFold(pokemonName, expectedPokemonName) {
		err := publishSuccessfulCatch(s, m, expectedPokemonName)
		if err != nil {
			log.Error(err)
		}

		log.Debug("Removing channel cache")
		h.channelCache.Del(m.ChannelID)
	} else {
		log.Info("TODO - block the same user from retrying catch")

		err := publichCatchFailure(s, m)
		if err != nil {
			log.Error(err)
		}
	}
}
