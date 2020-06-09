package spawn

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type pokemonSpawner struct {
	chance         float64
	channelCache   *ristretto.Cache
	debounceWindow time.Duration
}

func newPokemonSpawner(chance float64, debounceWindow time.Duration) (*pokemonSpawner, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e2,
		MaxCost:     1e2,
		BufferItems: 64,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create channel cache")
	}
	return &pokemonSpawner{
		chance:         chance,
		channelCache:   cache,
		debounceWindow: debounceWindow,
	}, nil
}

func (p *pokemonSpawner) spawn(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// check if channel id is still in cache, if it is
	// we are still in the debounce window
	_, exists := p.channelCache.Get(m.ChannelID)
	if exists {
		return false
	}

	// randomly determine if a pokemon is spawned
	if rand.Float64() > p.chance {
		return false
	}

	// spawn a pokemon!
	pokemonID := rand.Intn(890) + 1
	msg := discordgo.MessageEmbed{
		Title:       "‌‌A wild pokémon has appeared!",
		Description: "Guess the pokémon аnd type `p!ka catch <pokémon> with <ball>` to cаtch it!",
		Color:       0x008080,
		Image: &discordgo.MessageEmbedImage{
			URL: fmt.Sprintf("https://pokeres.bastionbot.org/images/pokemon/%d.png", pokemonID),
		},
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Error(err)
		return false
	}

	// TODO - add pokemon to either channel cache or object.

	// add channel id to cache, set to expire after the debounce window
	p.channelCache.SetWithTTL(m.ChannelID, struct{}{}, 1, p.debounceWindow)

	return true
}
