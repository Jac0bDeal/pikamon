package spawn

import (
	"fmt"
	"github.com/Jac0bDeal/pikamon/internal/pikamon/util"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	log "github.com/sirupsen/logrus"
)

type pokemonSpawner struct {
	chance               float64
	channelCache         *ristretto.Cache
	minimumSpawnDuration time.Duration
}

func newPokemonSpawner(botCache *util.BotCache, chance float64, minimumSpawnDuration time.Duration) (*pokemonSpawner, error) {
	return &pokemonSpawner{
		chance:               chance,
		channelCache:         botCache.ChannelCache,
		minimumSpawnDuration: minimumSpawnDuration,
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
	pokemonID := rand.Intn(807) + 1
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

	// add channel id to cache, set to expire after the debounce window
	log.Infof("Adding pokemon with id %d to channel cache", pokemonID)
	p.channelCache.SetWithTTL(m.ChannelID, pokemonID, 1, p.minimumSpawnDuration)

	return true
}
