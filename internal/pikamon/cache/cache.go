package cache

import (
	"github.com/Jac0bDeal/pikamon/internal/pikamon/config"
	"github.com/pkg/errors"

	"github.com/dgraph-io/ristretto"
)

// Cache holds the various caches for the Bot.
type Cache struct {
	Channel *ristretto.Cache
}

// New builds and returns a pointer to a Cache object.
func New(cfg *config.Config) (*Cache, error) {
	// Create our bot cache for channels
	channelCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.Cache.Channel.NumCounters,
		MaxCost:     cfg.Cache.Channel.MaxCost,
		BufferItems: cfg.Cache.Channel.BufferItems,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to build channel cache")
	}

	return &Cache{
		Channel: channelCache,
	}, nil
}

// Close closes all of the caches contained in Cache
func (c *Cache) Close() {
	c.Channel.Close()
}
