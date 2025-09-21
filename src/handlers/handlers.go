package handlers

import (
	"slices"

	"cox/src/constants"
	"cox/src/parsers"

	"github.com/bwmarrin/discordgo"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	var p *parsers.MessageParser

	p = parsers.NewMessageParser(s, m)
	p.Handle()
}

func MessageReactionAddHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Ignore reactions from the bot itself
	if r.UserID == s.State.User.ID {
		return
	}

	if !slices.Contains(constants.ASSIGN_ROLE_MESSAGE_IDS, r.MessageID) {
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

	if !slices.Contains(constants.ASSIGN_ROLE_MESSAGE_IDS, r.MessageID) {
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

func AccountMessageCreateHandler(pubChannel chan<- *discordgo.MessageCreate) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore messages from the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !slices.Contains(constants.SECOND_BOT_LISTENING_CHANNEL_IDS, m.ChannelID) {
			return
		}

		pubChannel <- m
	}
}

func BufferedMessageCreateHandler(subChannel <-chan *discordgo.MessageCreate, s *discordgo.Session) {
	go func() {
		for msg := range subChannel {
			p := parsers.NewMessageParser(s, msg)
			p.Handle()
		}
	}()
}
