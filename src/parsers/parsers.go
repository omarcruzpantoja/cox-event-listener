package parsers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	CoxCommand = "/cox-listener"

	RoleIdRegex = regexp.MustCompile(`<@&\d+>`)

	roleSetupRegex   = regexp.MustCompile("/cox-listener role-setup")
	titleRegex       = regexp.MustCompile(`title="([^"]+)"`)
	descriptionRegex = regexp.MustCompile(`description="([^"]+)"`)
	roleEmojiRegex   = regexp.MustCompile(`emoji\d+="([^"]+)"`)
	roleOptionRegex  = regexp.MustCompile(`option\d+="([^"]+)"`)
	roleRegex        = regexp.MustCompile(`role\d+="([^"]+)"`)
	getPositionId    = regexp.MustCompile(`.*?(\d+)$`)

	helpRegex = regexp.MustCompile("/cox-listener help")
)

type MessageParser struct {
	s *discordgo.Session
	m *discordgo.MessageCreate
}

func NewMessageParser(session *discordgo.Session, message *discordgo.MessageCreate) *MessageParser {
	return &MessageParser{
		s: session,
		m: message,
	}
}

func (mp *MessageParser) Handle() {
	if roleSetupRegex.Match([]byte(mp.m.Content)) {
		mp.roleSetupHandler()
	} else if helpRegex.Match([]byte(mp.m.Content)) {
		mp.helpHandler()
	}
}

func (mp *MessageParser) helpHandler() {
	commands := `--
**Available commands**:
1. **/cox-listener help** 			 - Prints available commands and description
2. **/cox-listener role-setup title="" description="" option1="" emoji1 ...** - Creates a message allowing users to receive role by selecting emojis.
																																	 Use add role descriptions and its corresponding emoji using optionX="" and emojiX=""
																																	 Example. option1="10x drop rates" emoji="💸"
`

	_, err := mp.s.ChannelMessageSend(mp.m.ChannelID, commands)

	if err != nil {
		log.Printf("Failure sending message (helpHandler): %s", err)
	}

}

func (mp *MessageParser) roleSetupHandler() {
	/*
		Full message example:

		/cox-listener role-setup
		title="Choose events you'd like to be notified on"
		description="Use the reaction below if you want to get notifications every time there is an X event in the game."

		option1="All money rates events" emoji1="💸"
		option2="X10+ money rates" emoji2="💰"
		option3="All drop rate events" emoji3="🏹"
		option4="All 10x+ drop rate events" emoji4="💎"
	*/

	var (
		splitRes       []string
		msgDescription []string
		objId          []string
		id             int
		err            error
	)

	title := titleRegex.FindString(mp.m.Content)
	description := descriptionRegex.FindString(mp.m.Content)
	emojis := roleEmojiRegex.FindAllString(mp.m.Content, -1)
	options := roleOptionRegex.FindAllString(mp.m.Content, -1)
	foundRoles := roleRegex.FindAllString(mp.m.Content, -1)

	reactions := make([]string, len(emojis))
	roleDescriptions := make([]string, len(options))
	roles := make([]string, len(options))

	/*
		---- Parse emojis -----
		Example:
		emoji1="💸"
		emoji2="💰"
		emoji3="🏹"
		emoji4="💎"
	*/
	for _, emoji := range emojis {
		splitRes = strings.Split(emoji, "=")
		objId = getPositionId.FindStringSubmatch(splitRes[0])

		if len(objId) > 1 {
			id, err = strconv.Atoi(objId[1])
			if err != nil {
				log.Printf("Unable to find option id (roleSetupHandler): %s\n", err)
				return
			}
		}

		reactions[id-1] = strings.Trim(splitRes[1], `"`)
	}

	/*
		---- Role description options -----
		Example:
		option1="All money rates events"
		option2="X10+ money rates"
		option3="All drop rate events"
		option4="All 10x+ drop rate events"
	*/
	for _, opt := range options {
		splitRes = strings.Split(opt, "=")
		objId = getPositionId.FindStringSubmatch(splitRes[0])

		if len(objId) > 1 {
			id, err = strconv.Atoi(objId[1])
			if err != nil {
				log.Printf("Unable to find option id (roleSetupHandler): %s\n", err)
				return
			}
		}

		roleDescriptions[id-1] = strings.Trim(splitRes[1], `"`)
	}

	/*
		---- Role options -----
		Example:
		role1="MoneyRate All"
		role2="MoneyRate 5x+"
		role3="DropRate All"
		role4="MoneyRate 10x+"
	*/
	for _, r := range foundRoles {
		splitRes = strings.Split(r, "=")
		objId = getPositionId.FindStringSubmatch(splitRes[0])

		if len(objId) > 1 {
			id, err = strconv.Atoi(objId[1])
			if err != nil {
				log.Printf("Unable to find role id (roleSetupHandler): %s\n", err)
				return
			}
		}

		roles[id-1] = strings.Trim(splitRes[1], `"`)
	}

	// ---- Get all roles ----
	guild, err := mp.s.State.Guild(mp.m.GuildID)

	if err != nil {
		log.Printf("Unable to get guild data (roleSetupHandler)")
		return
	}

	for _, role := range guild.Roles {
		for index := range len(roles) {
			if roles[index] == role.Name {
				roles[index] = fmt.Sprintf("<@&%s>", role.ID)
			}
		}
	}

	// ---- Create discord message -----
	msgDescription = append(msgDescription, argParser(title))
	msgDescription = append(msgDescription, fmt.Sprintf(">>> **%s**", argParser(title)))
	msgDescription = append(msgDescription, fmt.Sprintf("%s\n", argParser(description)))

	for i := range len(roleDescriptions) {
		msgDescription = append(
			msgDescription,
			fmt.Sprintf(("%s: **%s** (%s)"), reactions[i], roleDescriptions[i], roles[i]),
		)
	}

	// Send a message including message description and role descriptions (and its emoji)
	msg, err := mp.s.ChannelMessageSend(mp.m.ChannelID, strings.Join(msgDescription, "\n"))
	if err != nil {
		log.Printf("Failure sending message (roleSetupHandler): %s\n", err)
		return
	}

	// Add reactions to that message
	for _, r := range reactions {
		err := mp.s.MessageReactionAdd(mp.m.ChannelID, msg.ID, r)
		if err != nil {
			log.Printf("Failure reacting to message (roleSetupHandler): %s\n", err)
		}
	}

	// Delete command message
	err = mp.s.ChannelMessageDelete(mp.m.ChannelID, mp.m.ID)
	if err != nil {
		log.Printf("Failure deleting message (roleSetupHandler): %s\n", err)
	}

}

func argParser(s string) string {
	splitRes := strings.Split(s, "=")
	return strings.Trim(splitRes[1], `"`)
}
