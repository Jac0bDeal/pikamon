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

func pokemonExpiredFailure(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	expireMessageIndex := rand.Intn(len(pokemonExpiredMessages))
	expireMessage := "You have failed to catch the pokemon! " + pokemonExpiredMessages[expireMessageIndex]
	msg := discordgo.MessageEmbed{
		Title:       "The Pokemon has run away!",
		Description: expireMessage,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Error(err)
		return false
	}
	return true
}

func catchFailure(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	catchFailureMessageIndex := rand.Intn(len(catchFailureMessages))
	expireMessage := fmt.Sprintf("%s has failed to catch the pokemon! %s", m.Author.Username, catchFailureMessages[catchFailureMessageIndex])
	msg := discordgo.MessageEmbed{
		Description: expireMessage,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Error(err)
		return false
	}
	return true
}

func publishSuccessfulCatch(s *discordgo.Session, m *discordgo.MessageCreate, pName string) bool {
	// TODO - save to database
	catchMessage := fmt.Sprintf("Congratulations %s! You caught a %s", m.Author.Username, pName)
	msg := discordgo.MessageEmbed{
		Description: catchMessage,
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Error(err)
		return false
	}
	return true
}

func catch(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check if channel id is still in cache. If it is not then there is nothing to catch
	p, exists := util.BotMetadata.ChannelCache.Get(m.ChannelID)
	if !exists {
		publishFailure := pokemonExpiredFailure(s, m)
		if !publishFailure {
			log.Error("Unable to publish pokemonExpiredFailure message to channel")
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

	// The pokemon name specified by the person trying to catch it
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
		publishFailure := publishSuccessfulCatch(s, m, expectedPokemonName)
		if !publishFailure {
			log.Error("Unable to publish publishSuccessfulCatch message to channel")
		}

		log.Debug("Removing channel cache")
		util.BotMetadata.ChannelCache.Del(m.ChannelID)
	} else {
		log.Info("TODO - block the same user from retrying catch")

		publishFailure := catchFailure(s, m)
		if !publishFailure {
			log.Error("Unable to publish catchFailure message to channel")
		}
	}
}
