package handlers

import (
	"cox/src/parsers"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	var p *parsers.MessageParser

	// fmt.Printf("[%s] %s: %s...\n", m.ChannelID, m.Author.Username, m.Content)

	// First check if the command is a cox-listener, if so process it
	if strings.HasPrefix(m.Content, parsers.CoxCommand) {
		p = parsers.NewMessageParser(s, m)
		p.Handle()
		return
	}

	// Alternatively, only process messages from the specified channels
	if !slices.Contains(LISTENING_CHANNEL_IDS, m.ChannelID) {
		// If the reaction to the message is not included in the expected message ids, don't do anything
		return
	}
}

func MessageReactionAddHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Ignore reactions from the bot itself
	if r.UserID == s.State.User.ID {
		return
	}

	if !slices.Contains(ASSIGN_ROLE_MESSAGE_IDS, r.MessageID) {
		// If the reaction to the message is not included in the expected message ids, don't do anything
		return
	}

	performMessageReactionChange(s, MessageReactionChange{
		ChannelID: r.ChannelID,
		EmojiName: r.Emoji.Name,
		EmojiID:   r.Emoji.ID,
		GuildID:   r.GuildID,
		MessageID: r.MessageID,
		UserID:    r.UserID,

		add: true,
	})
}

func MessageReactionRemoveHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	// Ignore reactions from the bot itself
	if r.UserID == s.State.User.ID {
		return
	}

	if !slices.Contains(ASSIGN_ROLE_MESSAGE_IDS, r.MessageID) {
		// If the reaction to the message is not included in the expected message ids, don't do anything
		return
	}

	performMessageReactionChange(s, MessageReactionChange{
		ChannelID: r.ChannelID,
		EmojiID:   r.Emoji.ID,
		EmojiName: r.Emoji.Name,
		GuildID:   r.GuildID,
		MessageID: r.MessageID,
		UserID:    r.UserID,

		add: false,
	})
}
