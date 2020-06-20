package util

import (
	"github.com/dgraph-io/ristretto"
)

// Pikamon bot keyword
var CommandKeyword = "p!ka"

// Create a cache object for the Bot. May contain different caches of varying sizes (used for different purposes)
type BotCache struct {
	ChannelCache *ristretto.Cache
}

// Global metadata variable
var BotMetadata *BotCache
