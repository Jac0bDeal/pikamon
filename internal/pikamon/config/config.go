package config

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config contains the various values that are configurable for the bot.
type Config struct {
	Bot struct {
		MaximumSpawnDuration time.Duration
		SpawnChance          float64
		MaxPokemonID         int
	}
	Cache struct {
		Channel struct {
			NumCounters int64
			MaxCost     int64
			BufferItems int64
		}
	}
	Discord struct {
		Token string
	}
	Logging struct {
		Level string
	}
	Store struct {
		Type   string
		Sqlite struct {
			Location string
		}
	}
}

// GetConfig reads the config file and flags, then applies environment variable overrides.
func GetConfig() (*Config, error) {
	cfg := &Config{}

	// initialize config variables
	viper.SetEnvPrefix("PIKAMON")
	viper.SetConfigName("pikamon")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/pikamon")
	viper.AddConfigPath("./configs")

	// read config file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	cfg.Bot.MaximumSpawnDuration = viper.GetDuration("pikamon.bot.maximum-spawn-duration")
	cfg.Bot.SpawnChance = viper.GetFloat64("pikamon.bot.spawn-chance")
	cfg.Bot.MaxPokemonID = viper.GetInt("pikamon.bot.max-pokemon-id")

	cfg.Cache.Channel.NumCounters = viper.GetInt64("pikamon.cache.channel.number-counters")
	cfg.Cache.Channel.MaxCost = viper.GetInt64("pikamon.cache.channel.max-cost")
	cfg.Cache.Channel.BufferItems = viper.GetInt64("pikamon.cache.channel.buffer-size")

	cfg.Logging.Level = viper.GetString("pikamon.logging.level")

	cfg.Store.Type = viper.GetString("pikamon.store.type")
	cfg.Store.Sqlite.Location = viper.GetString("pikamon.store.sqlite.location")

	// define flags
	pflag.StringVarP(&cfg.Discord.Token, "token", "t", "", "Bot Token")

	// parse and bind flags
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, errors.Wrap(err, "failed to bind command line flags")
	}

	return cfg, nil
}
