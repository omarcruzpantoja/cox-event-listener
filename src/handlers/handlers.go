package handlers

import (
	"cox/src/parsers"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	var p *parsers.MessageParser

	fmt.Printf("[%s] %s: %s...\n", m.ChannelID, m.Author.Username, m.Content)

	if strings.HasPrefix(m.Content, parsers.CoxCommand) {
		p = parsers.NewMessageParser(s, m)
		p.Handle()
	}

	// Only process messages from the specified channel
	// if m.ChannelID == ChannelID {
	// }
}

func MessageReactionAddHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Ignore reactions from the bot itself
	if r.UserID == s.State.User.ID {
		return
	}

	fmt.Printf("Reaction added by %s to message %s in channel %s - roleid(%s) \n", r.UserID, r.MessageID, r.ChannelID, r.Emoji.ID)
	// Additional logic for handling reactions can be added here
}
