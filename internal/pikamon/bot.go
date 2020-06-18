package pikamon

import (
	"fmt"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/commands"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/spawn"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/util"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
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

	// Create bot cache
	createBotCache(cfg)

	// register discord handlers
	spawnListener, err := spawn.NewHandler(util.BotMetadata, cfg.Bot.SpawnChance, cfg.Bot.DebounceWindow)
	if err != nil {
		return nil, err
	}

	discord.AddHandler(commands.Handle)
	discord.AddHandler(spawnListener.Handle)

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

// Initialize bot cache object
func createBotCache(cfg *Config) {
	// Create our bot cache for channels
	channelCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.ChannelCache.NumCounters,
		MaxCost:     cfg.ChannelCache.MaxCost,
		BufferItems: cfg.ChannelCache.BufferItems,
	})
	if err != nil {
		log.Fatal(err, "failed to create channel cache")
	}

	util.BotMetadata = &util.BotCache{
		ChannelCache: channelCache,
	}
}
