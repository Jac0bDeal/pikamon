package pikamon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/cache"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/commands"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/config"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/spawn"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/store"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Bot listens to Discord and performs the various actions of Pikamon.
type Bot struct {
	cache   *cache.Cache
	store   store.Store
	discord *discordgo.Session
}

// New configures a Bot from the passed config, and returns it.
func New(cfg *config.Config) (*Bot, error) {
	log.Info("Creating Discord session...")
	authStr := fmt.Sprintf("Bot %s", cfg.Discord.Token)
	discord, err := discordgo.New(authStr)
	if err != nil {
		return nil, err
	}
	log.Info("Discord session created.")

	log.Info("Creating bot cache...")
	botCache, err := cache.New(cfg)
	if err != nil {
		return nil, err
	}
	log.Info("Bot cache created.")

	log.Info("Creating bot store...")
	botStore, err := store.New(cfg)
	if err != nil {
		return nil, err
	}
	log.Info("Bot store created.")

	log.Info("Registering bot handlers...")
	log.Debug("Registering commands handler...")
	commandsHandler := commands.NewHandler(
		botCache,
		botStore,
	)
	discord.AddHandler(commandsHandler.Handle)
	log.Debug("Commands handler registered.")
	log.Debug("Registering spawn handler...")
	spawnHandler := spawn.NewHandler(
		botCache,
		botStore,
		cfg.Bot.SpawnChance,
		cfg.Bot.MaximumSpawnDuration,
		cfg.Bot.MaxPokemonID,
	)
	discord.AddHandler(spawnHandler.Handle)
	log.Debug("Spawn handler registered.")

	return &Bot{
		cache:   botCache,
		discord: discord,
		store:   botStore,
	}, nil
}

// Run starts the bot, listens for a halt signal, and shuts down when the halt is received.
func (b *Bot) Run() error {
	if err := b.Start(); err != nil {
		return errors.Wrap(err, "failed to start bot")
	}

	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Info("Received stop signal, shutting down...")
	b.Stop()
	return nil
}

// Start opens the connection to the discord web socket.
func (b *Bot) Start() error {
	log.Info("Starting bot...")
	if err := b.store.Open(); err != nil {
		return errors.Wrap(err, "failed to open pikamon store")
	}

	log.Info("Opening connection to Discord...")
	if err := b.discord.Open(); err != nil {
		return errors.Wrap(err, "failed to open web socket connection to Discord")
	}
	log.Info("Connection to Discord established.")
	return nil
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() {
	log.Info("Stopping bot...")
	log.Info("Closing connection to Discord...")
	err := b.discord.Close()
	if err != nil {
		log.Error("Error closing discord api session: %v", err)
	}
	log.Info("Connection to Discord closed...")

	b.cache.Close()

	err = b.store.Close()
	if err != nil {
		log.Error("Error closing store session: %v", err)
	}
}
