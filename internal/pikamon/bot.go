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
	authStr := fmt.Sprintf("Bot %s", cfg.Discord.Token)
	discord, err := discordgo.New(authStr)
	if err != nil {
		return nil, err
	}

	// create bot cache
	botCache, err := cache.New(cfg)
	if err != nil {
		return nil, err
	}

	// create bot store
	botStore, err := store.New(cfg)
	if err != nil {
		return nil, err
	}

	// register discord handlers
	commandsHandler := commands.NewHandler(botCache.Channel)

	spawnHandler, err := spawn.NewHandler(
		botCache.Channel,
		cfg.Bot.SpawnChance,
		cfg.Bot.MaximumSpawnDuration,
		cfg.Bot.MaxPokemonID,
	)
	if err != nil {
		return nil, err
	}

	discord.AddHandler(commandsHandler.Handle)
	discord.AddHandler(spawnHandler.Handle)

	return &Bot{
		cache:   botCache,
		discord: discord,
		store:   botStore,
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
	b.Stop()
	return nil
}

// Start opens the connection to the discord web socket.
func (b *Bot) Start() error {
	if err := b.store.Open(); err != nil {
		return errors.Wrap(err, "failed to open pikamon store")
	}
	if err := b.discord.Open(); err != nil {
		return errors.Wrap(err, "failed to open web socket connection to Discord")
	}
	return nil
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() {
	err := b.discord.Close()
	if err != nil {
		log.Error("Error closing discord api session: %v", err)
	}

	b.cache.Close()

	err = b.store.Close()
	if err != nil {
		log.Error("Error closing store session: %v", err)
	}
}
