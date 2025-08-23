package handlers

import (
	"cox/src/parsers"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type MessageReactionChange struct {
	ChannelID string
	EmojiID   string
	EmojiName string
	GuildID   string
	MessageID string
	UserID    string

	add bool
}

func performMessageReactionChange(s *discordgo.Session, r MessageReactionChange) {
	var reactionRole string = ""

	// Get message to find the emoji-to-role mapping
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)

	if err != nil {
		log.Printf("Error retrieving message (MessageReaction%v): %v", r.add, err)
		return
	}

	// separate lines
	lines := strings.Split(message.Content, "\n")

	// find the line that contains the role (if any)
	for _, line := range lines {
		if strings.Contains(line, r.EmojiName) {
			reactionRole = parsers.RoleIdRegex.FindString(line)
			break
		}
	}

	guild, err := s.State.Guild(r.GuildID)

	if err != nil {
		// if fetching guild data errors revert the reaction change
		log.Printf("Error retrieving guild data (MessageReaction%v): %v", r.add, err)
		return
	}

	for _, role := range guild.Roles {

		if strings.Contains(reactionRole, role.ID) {

			if r.add {
				err = s.GuildMemberRoleAdd(r.GuildID, r.UserID, role.ID)

				if err != nil {
					log.Printf("Error giving role to user (MessageReaction%v): %v", r.add, err)
				}

			} else {
				err = s.GuildMemberRoleRemove(r.GuildID, r.UserID, role.ID)

				if err != nil {
					log.Printf("Error removing role to user (MessageReaction%v): %v", r.add, err)
				}
			}

			break
		}
	}
}
