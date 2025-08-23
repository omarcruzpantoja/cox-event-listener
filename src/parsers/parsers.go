package parsers

import (
	"cox/src/constants"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	CoxCommand = "/cox-listener"

	RoleIdRegex = regexp.MustCompile(`<@&\d+>`)

	roleSetupRegex     = regexp.MustCompile("/cox-listener role-setup")
	titleRegex         = regexp.MustCompile(`title="([^"]+)"`)
	descriptionRegex   = regexp.MustCompile(`description="([^"]+)"`)
	roleEmojiRegex     = regexp.MustCompile(`emoji\d+="([^"]+)"`)
	roleOptionRegex    = regexp.MustCompile(`option\d+="([^"]+)"`)
	roleRegex          = regexp.MustCompile(`role\d+="([^"]+)"`)
	getPositionIdRegex = regexp.MustCompile(`.*?(\d+)$`)

	helpRegex = regexp.MustCompile("/cox-listener help")

	// RATES MATCHERS
	getRateValueRegex       = regexp.MustCompile(`x\d+`)
	dropRateRegex           = regexp.MustCompile(`x\d+ drop rate`)
	goldMultiplierRateRegex = regexp.MustCompile(`x\d+ gold multiplier rate`)
	dbSocRateRegex          = regexp.MustCompile(`x\d+ dragonball soc rate`)
	metSocRateRegex         = regexp.MustCompile(`x\d+ meteor soc rate`)
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
	// First check for cox-listener commands
	if roleSetupRegex.MatchString(mp.m.Content) {
		mp.roleSetupHandler()
	} else if helpRegex.MatchString(mp.m.Content) {
		mp.helpHandler()

	} else {

		// If the reaction to the message is not included in the expected message ids, don't do anything
		// Can also add user id check to further scope down
		if !slices.Contains(constants.LISTENING_CHANNEL_IDS, mp.m.ChannelID) {
			return
		}

		if dropRateRegex.MatchString(mp.m.Content) {
			mp.dropRateHandler()
		} else if goldMultiplierRateRegex.MatchString(mp.m.Content) {
			mp.goldMultiplierRateHandler()
		} else if dbSocRateRegex.MatchString(mp.m.Content) {
			mp.dbSocRateHandler()
		} else if metSocRateRegex.MatchString(mp.m.Content) {
			mp.metSocRateHandler()
		}
	}

	// Example mesages
	// | x3 gold multiplier rate has started for 9 minutes!
	// | x21 drop rate has started for 13 minutes!
	// | x4 dragonball soc rate finished!
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
		objId = getPositionIdRegex.FindStringSubmatch(splitRes[0])

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
		objId = getPositionIdRegex.FindStringSubmatch(splitRes[0])

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
		objId = getPositionIdRegex.FindStringSubmatch(splitRes[0])

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
	msgDescription = append(msgDescription, argParser(title, "="))
	msgDescription = append(msgDescription, fmt.Sprintf(">>> **%s**", argParser(title, "=")))
	msgDescription = append(msgDescription, fmt.Sprintf("%s\n", argParser(description, "=")))

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

func (mp *MessageParser) dropRateHandler() {
	roles := mp.rolesGenerator("drop")

	for _, channelId := range constants.RELAY_MESSAGE_CHANNEL_IDS {

		_, err := mp.s.ChannelMessageSend(channelId, fmt.Sprintf("%s %s", strings.Join(roles, " "), mp.m.Content))

		if err != nil {
			log.Printf("Error sending message (DropRate Message) to channel (%s)\n", channelId)
		}
	}
}

func (mp *MessageParser) goldMultiplierRateHandler() {
	roles := mp.rolesGenerator("gold")

	for _, channelId := range constants.RELAY_MESSAGE_CHANNEL_IDS {

		_, err := mp.s.ChannelMessageSend(channelId, fmt.Sprintf("%s %s", strings.Join(roles, " "), mp.m.Content))

		if err != nil {
			log.Printf("Error sending message (GoldMultiplierRate Message) to channel (%s)\n", channelId)
		}
	}
}

func (mp *MessageParser) dbSocRateHandler() {
	roles := mp.rolesGenerator("dbSoc")

	for _, channelId := range constants.RELAY_MESSAGE_CHANNEL_IDS {

		_, err := mp.s.ChannelMessageSend(channelId, fmt.Sprintf("%s %s", strings.Join(roles, " "), mp.m.Content))

		if err != nil {
			log.Printf("Error sending message (DbSocRate Message) to channel (%s)\n", channelId)
		}
	}
}

func (mp *MessageParser) metSocRateHandler() {
	roles := mp.rolesGenerator("metSoc")

	for _, channelId := range constants.RELAY_MESSAGE_CHANNEL_IDS {

		_, err := mp.s.ChannelMessageSend(channelId, fmt.Sprintf("%s %s", strings.Join(roles, " "), mp.m.Content))

		if err != nil {
			log.Printf("Error sending message (MetSocRate Message) to channel (%s)\n", channelId)
		}
	}
}

func (mp *MessageParser) rolesGenerator(rateType string) []string {

	var (
		roles   []string
		allIds  []string
		tenXIds []string
	)

	if strings.Contains(mp.m.Content, "finished") {
		// if the message has "finished" in it, don't all notifications (roles)
		return roles
	}

	stringRate := argParser(getRateValueRegex.FindString(mp.m.Content), "x")
	rate, err := strconv.Atoi(stringRate)

	switch rateType {
	case "drop":
		allIds = constants.DROP_RATE_ROLE_IDS
		tenXIds = constants.DROP_RATE_10X_ROLE_IDS
	case "gold":
		allIds = constants.GOLD_MULTIPLIER_RATE_ROLE_IDS
		tenXIds = constants.GOLD_MULTIPLIER_RATE_5X_ROLE_IDS
	case "dbSoc":
		allIds = constants.DB_SOC_RATE_ROLE_IDS
		tenXIds = constants.DB_SOC_RATE_10X_ROLE_IDS
	case "metSoc":
		allIds = constants.MET_SOC_RATE_ROLE_IDS
		tenXIds = constants.MET_SOC_RATE_10X_ROLE_IDS
	}

	if err != nil {
		log.Printf("Error parsing drop rate value (%s).\n", stringRate)
		return roles
	}

	for _, allRateIds := range allIds {
		roles = append(roles, fmt.Sprintf("<@&%s>", allRateIds))
	}

	if rate >= 10 {
		for _, tenXRateIds := range tenXIds {
			roles = append(roles, fmt.Sprintf("<@&%s>", tenXRateIds))
		}
	}

	return roles
}

func argParser(s string, separator string) string {
	splitRes := strings.Split(s, separator)
	// trim is for cleanup
	return strings.Trim(splitRes[1], `"`)
}
