package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	authToken string
)

func init() {
	viper.SetEnvPrefix("pikamon")
	viper.SetConfigName("pikamon")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/pikamon")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	pflag.StringVarP(&authToken, "token", "t", "", "Bot Token")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatalf("Error binding command line flags: %s", err)
	}

	authToken = viper.GetString("token")
}

func main() {
	authStr := fmt.Sprintf("Bot %s", authToken)
	discord, err := discordgo.New(authStr)
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(messageCreate)

	if err = discord.Open(); err != nil {
		log.Fatalf("Error opening connection to Discord: %s", err)
		return
	}

	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Info("Received stop signal, shutting down...")

	if err := discord.Close(); err != nil {
		log.Fatalf("Error closing connection to Discord: %s", err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Debugf("Received message create: %+v", m.Message)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := strings.ToLower(m.Content)
	if content == "hi" || content == "hello" || content == "hey" {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Hey there!"); err != nil {
			log.Error(err)
		}
	}
}
