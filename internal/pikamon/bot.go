package pikamon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/appearance"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/commands"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Bot listens to Discord and performs the various actions of Pikamon.
type Bot struct {
	discord *discordgo.Session
}

// NewListener configures a Bot from the passed config, and returns it.
func New(cfg *Config) (*Bot, error) {
	authStr := fmt.Sprintf("Bot %s", cfg.Discord.Token)
	discord, err := discordgo.New(authStr)
	if err != nil {
		return nil, err
	}

	listener, err := appearance.NewListener(cfg.Bot.SpawnChance, cfg.Bot.DebounceWindow)
	if err != nil {
		return nil, err
	}

	// register discord handlers
	discord.AddHandler(commands.Handle)
	discord.AddHandler(listener.Handle)

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
