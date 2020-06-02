package main

import (
	"github.com/Jac0bDeal/pikamon/internal/logging"
	"github.com/Jac0bDeal/pikamon/internal/pikamon"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := pikamon.GetConfig()
	if err != nil {
		log.Fatal("Error getting Pikamon config: ", err)
	}

	if err := logging.Configure(cfg.Logging.Level); err != nil {
		log.Fatal("Error configuring logger")
	}

	bot, err := pikamon.New(cfg)
	if err != nil {
		log.Fatal("Error configuring Pikamon bot: ", err)
	}

	if err := bot.Run(); err != nil {
		log.Fatal("Bot encountered fatal error: ", err)
	}
}
