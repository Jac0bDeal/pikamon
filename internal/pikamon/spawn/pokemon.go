package spawn

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/constants"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	log "github.com/sirupsen/logrus"
)

type pokemonSpawner struct {
	chance               float64
	channelCache         *ristretto.Cache
	maximumSpawnDuration time.Duration
	maxPokemonID         int
}

func newPokemonSpawner(
	channelCache *ristretto.Cache,
	chance float64,
	maximumSpawnDuration time.Duration,
	maxPokemonID int,
) *pokemonSpawner {
	return &pokemonSpawner{
		chance:               chance,
		channelCache:         channelCache,
		maximumSpawnDuration: maximumSpawnDuration,
		maxPokemonID:         maxPokemonID,
	}
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
	pokemonID := rand.Intn(p.maxPokemonID) + 1
	msg := discordgo.MessageEmbed{
		Title:       "‌‌A wild pokémon has appeared!",
		Description: "Guess the pokémon аnd type `p!ka catch <pokémon>` to cаtch it!",
		Color:       constants.MessageColor,
		Image: &discordgo.MessageEmbedImage{
			URL: fmt.Sprintf("https://pokeres.bastionbot.org/images/pokemon/%d.png", pokemonID),
		},
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, &msg); err != nil {
		log.Error(err)
		return false
	}

	// add channel id to cache, set to expire after the debounce window
	log.Debugf("Adding pokemon with id %d to channel cache", pokemonID)
	p.channelCache.SetWithTTL(m.ChannelID, pokemonID, 1, p.maximumSpawnDuration)

	return true
}
