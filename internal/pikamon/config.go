package pikamon

import (
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

	// define flags
	pflag.StringVarP(&cfg.Discord.Token, "token", "t", "", "Bot Token")

	// parse and bind flags
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, errors.Wrap(err, "failed to bind command line flags")
	}

	return cfg, nil
}
