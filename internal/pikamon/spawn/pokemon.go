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
	chance         float64
	channelCache   *ristretto.Cache
	debounceWindow time.Duration
}

type PokemonInfo struct {
	pokemonName string
	pokemonId   int
}

func newPokemonSpawner(botCache *util.BotCache, chance float64, debounceWindow time.Duration) (*pokemonSpawner, error) {
	return &pokemonSpawner{
		chance:         chance,
		channelCache:   botCache.ChannelCache,
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
	// TODO - Use API to get name of pokemon with ID pokemonID
	pokemonName := "SOME_POKEMON_NAME"
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

	// create pokemon info object
	//var pokemonObj *PokemonInfo = &PokemonInfo{
	//	pokemonName: pokemonName,
	//	pokemonId:   pokemonID,
	//}

	log.Infof("Adding pokemon %s with id %d to cache", pokemonName, pokemonID)

	// add channel id to cache, set to expire after the debounce window
	p.channelCache.SetWithTTL(m.ChannelID, pokemonID, 1, p.debounceWindow)

	return true
}
