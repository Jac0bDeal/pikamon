package pikamon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/commands"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/spawn"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Bot listens to Discord and performs the various actions of Pikamon.
type Bot struct {
	discord *discordgo.Session
}

// New configures a Bot from the passed config, and returns it.
func New(cfg *Config) (*Bot, error) {
	authStr := fmt.Sprintf("Bot %s", cfg.Discord.Token)
	discord, err := discordgo.New(authStr)
	if err != nil {
		return nil, err
	}

	// TODO - We need to add the message ID of the spawned pokemon to a struct or other object to reference later
	// If we know the message ID we can use the following to get the posted pokemon so we know what we are catching - https://discord.com/developers/docs/resources/channel#get-channel-message
	// May also need to use the message metadata to determine if pokemon is still there unless that is what the DebounceWindow is
	// we likely also need a way to determine when to clean up the pokemon that was last spawned. Perhaps that is its own issue.

	listener, err := spawn.NewHandler(cfg.Bot.SpawnChance, cfg.Bot.DebounceWindow)
	if err != nil {
		return nil, err
	}

	// TODO - Create catcher handler

	// register discord handlers
	discord.AddHandler(commands.Handle)
	discord.AddHandler(listener.Handle)
	// TODO - Add handler for catcher

	return &Bot{
		discord: discord,
	}, nil
}

// Run starts the bot, listens for a halt signal, and shuts down when the halt is received.
func (b *Bot) Run() error {
	log.Info("Starting bot...")
	if err := b.Start(); err != nil {
		return errors.Wrap(err, "failed to start bot")
	}

	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Info("Received stop signal, shutting down...")
	if err := b.Stop(); err != nil {
		return errors.Wrap(err, "failed to stop bot gracefully")
	}
	return nil
}

// Start opens the connection to the discord web socket.
func (b *Bot) Start() error {
	if err := b.discord.Open(); err != nil {
		return errors.Wrap(err, "failed to open web socket connection to Discord")
	}
	return nil
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() error {
	return b.discord.Close()
}
