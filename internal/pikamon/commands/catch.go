package commands

import (
	"fmt"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/util"
	"github.com/bwmarrin/discordgo"
	"github.com/mtslzr/pokeapi-go"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
)

var pokemonExpiredMessages = []string{
	"The pokemon heard you coming and ran for the hills!",
	"The pokemon thought you were ugly and didn't want to be caught by an ugly person!",
	"The pokemon likes Physics (mistakenly) and thought you were too dumb to understand it so it ran!",
}

var catchFailureMessages = []string{
	"You may need to have your eyes checked!",
	"Better luck next time!",
	"Have you tried getting good?",
}

func publishExpiredPokemon(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	expireMessageIndex := rand.Intn(len(pokemonExpiredMessages))
	expireMessage := "You have failed to catch the pokemon! " + pokemonExpiredMessages[expireMessageIndex]
	msg := discordgo.MessageEmbed{
		Title:       "The Pokemon has run away!",
		Description: expireMessage,
		Color:       0x008080,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msg)
	return err
}

func publichCatchFailure(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	catchFailureMessageIndex := rand.Intn(len(catchFailureMessages))
	expireMessage := fmt.Sprintf("%s has failed to catch the pokemon! %s", m.Author.Username, catchFailureMessages[catchFailureMessageIndex])
	msg := discordgo.MessageEmbed{
		Description: expireMessage,
		Color:       0x008080,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msg)
	return err
}

func publishSuccessfulCatch(s *discordgo.Session, m *discordgo.MessageCreate, pName string) (err error) {
	// TODO - save to database
	catchMessage := fmt.Sprintf("Congratulations %s! You caught a %s!", m.Author.Username, strings.Title(pName))
	msg := discordgo.MessageEmbed{
		Description: catchMessage,
		Color:       0x008080,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msg)
	return err
}

func catch(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check if channel id is still in cache. If it is not then there is nothing to catch
	p, exists := util.BotMetadata.ChannelCache.Get(m.ChannelID)
	if !exists {
		err := publishExpiredPokemon(s, m)
		if err != nil {
			log.Error(err)
		}
		return
	}
	log.Debug("Pokemon cache still exists. Attempting catch.")

	// Create pokemon information object from the cache
	var pokemonId = p.(int)
	pInfo, err := pokeapi.Pokemon(strconv.Itoa(pokemonId))
	if err != nil {
		log.Error(err)
	}

	// Get everything after the "catch" command
	text := strings.TrimSpace(strings.ToLower(m.Content))
	commandText := strings.TrimSpace(text[len(util.CommandKeyword):])
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
	}).Info("Trying to catch a pokemon!")

	// Perform catch attempt
	var expectedPokemonName string = pInfo.Name
	if strings.EqualFold(pokemonName, expectedPokemonName) {
		err := publishSuccessfulCatch(s, m, expectedPokemonName)
		if err != nil {
			log.Error(err)
		}

		log.Debug("Removing channel cache")
		util.BotMetadata.ChannelCache.Del(m.ChannelID)
	} else {
		log.Info("TODO - block the same user from retrying catch")

		err := publichCatchFailure(s, m)
		if err != nil {
			log.Error(err)
		}
	}
}
