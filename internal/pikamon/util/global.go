package util

import (
	"github.com/dgraph-io/ristretto"
)

type PokemonInfo struct {
	PokemonName string
	PokemonId   int
}

// Create a cache object for the Bot. May contain different caches of varying sizes (used for different purposes)
type BotCache struct {
	ChannelCache *ristretto.Cache
	Sample       string
}

// Global metadata variable
var BotMetadata *BotCache

// TODO - make global
var CommandKeyword = "p!ka"
