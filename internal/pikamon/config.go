package pikamon

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config contains the various values that are configurable for the bot.
type Config struct {
	Discord struct {
		Token string
	}
	Logging struct {
		Level string
	}
	Bot struct {
		MaximumSpawnDuration time.Duration
		SpawnChance          float64
		MaxPokemonID         int
	}

	ChannelCache struct {
		NumCounters int64
		MaxCost     int64
		BufferItems int64
	}
}

// GetConfig reads the config file and flags, then applies environment variable overrides.
func GetConfig() (*Config, error) {
	cfg := &Config{}

	// initialize config variables
	viper.SetEnvPrefix("pikamon")
	viper.SetConfigName("pikamon")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/pikamon")
	viper.AddConfigPath("./configs")

	// read config file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}
	cfg.Logging.Level = viper.GetString("pikamon.logging.level")
	cfg.Bot.MaximumSpawnDuration = viper.GetDuration("pikamon.bot.maximum-spawn-duration")
	cfg.Bot.SpawnChance = viper.GetFloat64("pikamon.bot.spawn-chance")
	cfg.Bot.MaxPokemonID = viper.GetInt("pikamon.bot.max-pokemon-id")
	cfg.ChannelCache.NumCounters = viper.GetInt64("pikamon.channel-cache.number-counters")
	cfg.ChannelCache.MaxCost = viper.GetInt64("pikamon.channel-cache.max-cost")
	cfg.ChannelCache.BufferItems = viper.GetInt64("pikamon.channel-cache.buffer-size")

	// define flags
	pflag.StringVarP(&cfg.Discord.Token, "token", "t", "", "Bot Token")

	// parse and bind flags
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, errors.Wrap(err, "failed to bind command line flags")
	}

	return cfg, nil
}
