package appearance

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Listener listens to non-Pikamon messages in a channel, and randomly causes
// pokemon to appear in that channel when messages are sent.
type Listener struct {
	channelCache   *ristretto.Cache
	debounceWindow time.Duration
	spawnChance    float64
}

// NewListener constructs and returns a new Listener that makes pokemon appear in
// a channel pseudo-randomly.
func NewListener(spawnChance float64, debounceWindow time.Duration) (*Listener, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e2,
		MaxCost:     1e2,
		BufferItems: 64,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create channel cache")
	}
	return &Listener{
		channelCache:   cache,
		debounceWindow: debounceWindow,
		spawnChance:    spawnChance,
	}, nil
}

// Handle is the handler function registered on the discord bot that
// processes incoming messages.
func (l *Listener) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if channel id is still in cache
	_, exists := l.channelCache.Get(m.ChannelID)
	if exists {
		return
	}

	// randomly determine if a pokemon is spawned
	if rand.Float64() > l.spawnChance {
		return
	}

	// add channel id to cache, set to expire after the debounce window
	l.channelCache.SetWithTTL(m.ChannelID, struct{}{}, 1, l.debounceWindow)

	// spawn a pokemon!
	pokemonID := rand.Intn(964) + 1
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
	}
}
