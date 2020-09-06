package cache

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/config"

	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Cache holds the various caches for the Bot.
type Cache struct {
	Channel *ristretto.Cache
}

// New builds and returns a pointer to a Cache object.
func New(cfg *config.Config) (*Cache, error) {
	log.Debug("Creating channel cache...")
	channelCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.Cache.Channel.NumCounters,
		MaxCost:     cfg.Cache.Channel.MaxCost,
		BufferItems: cfg.Cache.Channel.BufferItems,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to build channel cache")
	}
	log.Debug("Channel cache created.")

	return &Cache{
		Channel: channelCache,
	}, nil
}

// Close closes all of the caches contained in Cache
func (c *Cache) Close() {
	log.Info("Closing cache...")
	c.Channel.Close()
	log.Info("Cache closed.")
}
