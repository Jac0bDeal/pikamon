package catch

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"

	"github.com/Jac0bDeal/pikamon/internal/pikamon/commands"
)

type catcher interface {
	catch(*discordgo.Session, *discordgo.MessageCreate) bool
}

// Handler listens to Pikamon catch messages in a channel and attempts
// the catch operation with the specified pokeball.
type Handler struct {
	catchers []catcher
}

// TODO - first just build the catch. have it identify a catch message and print back to the channel "You caught it"

// CatchHandler constructs and returns a new Handler that catches pokemon in the channels.
// TODO - figure out if this should take the channel cache from the spawner
//func CatchHandler(s *discordgo.Session, m *discordgo.MessageCreate) (*Handler, error) {
func CatchHandler(x int) (*Handler, error) {
	return &Handler{}, nil
}

// CatchHandle is the handler function registered on the discord bot that
// processes incoming messages and checks if we are catching a spawned pokemon
func (h *Handler) CatchHandle(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == sess.State.User.ID {
		return
	}

	//fmt.Printf("Author ID: %s\n", m.Author.ID)
	//fmt.Printf("Message ID: %s\n", m.ID)
	//fmt.Printf("Message content: %s\n", m.Content)

	text := strings.TrimSpace(strings.ToLower(m.Content))
	fmt.Printf("Message text: %s\n", text)

	// ignore all messages not prefixed with bot command keyword
	if !strings.HasPrefix(text, commands.CommandKeyword) {
		fmt.Printf("CatchHandle: No pika command specified")
		return
	}
	commandText := strings.TrimSpace(text[len(commands.CommandKeyword):])
	if commandText == "" {
		return
	}

	//fmt.Printf("CatchHandle: Command text = %s\n", commandText)
}
