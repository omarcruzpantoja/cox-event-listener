package parsers

import (
	"cox/src/constants"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (

	// COX COMMAND MATCHERS
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

	initializeMessageRegex = regexp.MustCompile("/cox-listener message-init")

	// RATES MATCHERS
	getRateValueRegex       = regexp.MustCompile(`x\d+`)
	dropRateRegex           = regexp.MustCompile(`x\d+ drop rate`)
	goldMultiplierRateRegex = regexp.MustCompile(`x\d+ gold multiplier rate`)
	dbSocRateRegex          = regexp.MustCompile(`x\d+ dragonball soc rate`)
	metSocRateRegex         = regexp.MustCompile(`x\d+ meteor soc rate`)

	// WAR MATCHERS
	cityWarWinRegex = regexp.MustCompile("(.*) won (AC|PC|BI|DC) City War")

	// LAB BOSS WATCHES
	labBossSpawnRegex = regexp.MustCompile("(Howler|Gibbon|Talon|NagaLord) lab boss has been spawned random.")
	labBossKillRegex  = regexp.MustCompile("(Howler|Gibbon|Talon|NagaLord) lab boss has been killed by (.*)")

	// TEXT VARS
	cwMessageStart = `-
**CW Results**
PC: 
AC: 
DC: 
BI: 

**Next City War**: 10:00 UTC`

	labBossStart = `-
Lab 1 (Gibbon):
Lab 2 (Naga):
Lab 3 (Talon):
Lab 4 (Howler):
`

	MINUTES_FORMAT = "15:04"
	HOURS_FORMAT   = "15:00"
)

type MessageParser struct {
	s *discordgo.Session
	m *discordgo.MessageCreate
}

type ServerRoleMap = map[string][]string

func NewMessageParser(session *discordgo.Session, message *discordgo.MessageCreate) *MessageParser {
	return &MessageParser{
		s: session,
		m: message,
	}
}

func (mp *MessageParser) Handle() {
	// First check for cox-listener commands
	if roleSetupRegex.MatchString(mp.m.Content) && mp.isAdminUser() {
		mp.roleSetupHandler()
	} else if helpRegex.MatchString(mp.m.Content) && mp.isAdminUser() {
		mp.helpHandler()
	} else if initializeMessageRegex.MatchString(mp.m.Content) {
		mp.initializeMessageHandler()
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
		} else if cityWarWinRegex.MatchString(mp.m.Content) {
			mp.cityWarHandler()
		} else if labBossSpawnRegex.MatchString(mp.m.Content) {
			mp.labBossHandler(false)
		} else if labBossKillRegex.MatchString((mp.m.Content)) {
			mp.labBossHandler(true)
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
3. **/cox-listener message-init** - makes the bot send a message in the channel. Use _cw-init_ to initialize a city war format message
																		Example: _/cox-listener message-init cw-init_
`

	_, err := mp.s.ChannelMessageSend(mp.m.ChannelID, commands)

	if err != nil {
		log.Printf("Failure sending message (helpHandler): %s.\n", err)
	}
}

func (mp *MessageParser) initializeMessageHandler() {
	if strings.Contains(mp.m.Content, "cw-init") {
		mp.s.ChannelMessageSend(mp.m.ChannelID, cwMessageStart)
	} else if strings.Contains(mp.m.Content, "lab-boss-start") {
		mp.s.ChannelMessageSend(mp.m.ChannelID, labBossStart)
	} else {
		mp.s.ChannelMessageSend(mp.m.ChannelID, "init")
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
				log.Printf("Unable to find option id (roleSetupHandler): %s.\n", err)
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
				log.Printf("Unable to find option id (roleSetupHandler): %s.\n", err)
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
				log.Printf("Unable to find role id (roleSetupHandler): %s.\n", err)
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
		log.Printf("Failure sending message (roleSetupHandler): %s.\n", err)
		return
	}

	// Add reactions to that message
	for _, r := range reactions {
		err := mp.s.MessageReactionAdd(mp.m.ChannelID, msg.ID, r)
		if err != nil {
			log.Printf("Failure reacting to message (roleSetupHandler): %s.\n", err)
		}
	}

	// Delete command message
	err = mp.s.ChannelMessageDelete(mp.m.ChannelID, mp.m.ID)
	if err != nil {
		log.Printf("Failure deleting message (roleSetupHandler): %s.\n", err)
	}

}

func (mp *MessageParser) dropRateHandler() {
	roles := mp.rolesGeneratorRates("drop")

	for _, channel := range constants.RELAY_MESSAGE_CHANNEL_IDS {

		_, err := mp.s.ChannelMessageSend(
			channel.ChannelID,
			fmt.Sprintf("%s %s", strings.Join(roles[channel.ServerID], " "), mp.m.Content),
		)

		if err != nil {
			log.Printf("Error sending message (DropRate Message) to channel (%s). err: %s\n", channel.ChannelID, err)
		}
	}
}

func (mp *MessageParser) goldMultiplierRateHandler() {
	roles := mp.rolesGeneratorRates("gold")

	for _, channel := range constants.RELAY_MESSAGE_CHANNEL_IDS {

		_, err := mp.s.ChannelMessageSend(
			channel.ChannelID, fmt.Sprintf("%s %s", strings.Join(roles[channel.ServerID], " "), mp.m.Content),
		)

		if err != nil {
			log.Printf(
				"Error sending message (GoldMultiplierRate Message) to channel (%s). err: %s\n", channel.ChannelID, err,
			)
		}
	}
}

func (mp *MessageParser) dbSocRateHandler() {
	roles := mp.rolesGeneratorRates("dbSoc")

	for _, channel := range constants.RELAY_MESSAGE_CHANNEL_IDS {

		_, err := mp.s.ChannelMessageSend(
			channel.ChannelID, fmt.Sprintf("%s %s", strings.Join(roles[channel.ServerID], " "), mp.m.Content),
		)

		if err != nil {
			log.Printf("Error sending message (DbSocRate Message) to channel (%s). err: %s\n", channel.ChannelID, err)
		}
	}
}

func (mp *MessageParser) metSocRateHandler() {
	roles := mp.rolesGeneratorRates("metSoc")

	for _, channel := range constants.RELAY_MESSAGE_CHANNEL_IDS {

		_, err := mp.s.ChannelMessageSend(
			channel.ChannelID, fmt.Sprintf("%s %s", strings.Join(roles[channel.ServerID], " "), mp.m.Content),
		)

		if err != nil {
			log.Printf("Error sending message (MetSocRate Message) to channel (%s). err: %s\n", channel.ChannelID, err)
		}
	}
}

func (mp *MessageParser) labBossHandler(killed bool) {
	var (
		boss        string
		killer      string
		groups      []string
		message     string
		lab         string
		serverRoles ServerRoleMap
	)

	serverRoles = mp.roleGenerator("labBoss", nil)
	killer = ""

	for _, lbID := range constants.LAB_BOSS_CHANNEL_AND_MESSAGE_ID {
		labMessage, err := mp.s.ChannelMessage(lbID.ChannelID, lbID.MessageID)
		lines := strings.Split(labMessage.Content, "\n")

		if killed {
			groups = labBossKillRegex.FindStringSubmatch(mp.m.Content)
			boss = groups[1]
			if len(groups) > 2 {
				killer = groups[2]
			}
		} else {
			groups = labBossSpawnRegex.FindStringSubmatch(mp.m.Content)
			boss = groups[1]
		}

		if killed {
			message = fmt.Sprintf(
				"DEAD respawns in: (%s)-(%s) KILLER: %s",
				discordTimestamp(time.Now(), 5, 0),
				discordTimestamp(time.Now(), 8, 0),
				killer,
			)
		} else {
			message = fmt.Sprintf(
				"SPAWNED (%s)", discordTimestamp(time.Now(), 0, 0),
			)
		}

		switch boss {
		case "Gibbon":
			mp.roleGenerator("gibbon", &serverRoles)
			lines[1] = fmt.Sprintf("Lab 1 (Gibbon):\t%s", message)
			lab = "1"
		case "NagaLord":
			mp.roleGenerator("nagaLord", &serverRoles)
			lines[2] = fmt.Sprintf("Lab 2 (Naga):\t%s", message)
			lab = "2"
		case "Talon":
			mp.roleGenerator("talon", &serverRoles)
			lines[3] = fmt.Sprintf("Lab 3 (Talon):\t%s", message)
			lab = "3"
		case "Howler":
			mp.roleGenerator("howler", &serverRoles)
			lines[4] = fmt.Sprintf("Lab 4 (Howler):\t%s", message)
			lab = "4"
		}

		_, err = mp.s.ChannelMessageEdit(lbID.ChannelID, lbID.MessageID, strings.Join(lines, "\n"))

		if err != nil {
			log.Printf(
				"Failed to update message from channel (%s) message (%s) (labBossHandler)- err: %s.\n",
				lbID.ChannelID, lbID.MessageID, err,
			)
		}
	}

	if !killed {
		for _, channel := range constants.RELAY_MESSAGE_CHANNEL_IDS {

			_, err := mp.s.ChannelMessageSend(
				channel.ChannelID,
				fmt.Sprintf(
					"SPAWN %s (%s): Lab%s -> %s", strings.Join(serverRoles[channel.ServerID], " "),
					discordTimestamp(time.Now(), 0, 0),
					lab,
					boss,
				),
			)

			if err != nil {
				log.Printf(
					"Error sending message (labBossHandler Message) to channel (%s). err: %s\n", channel.ChannelID, err,
				)
			}
		}
	} else {
		// If killed, also send who killed it
		for _, channel := range constants.RELAY_MESSAGE_CHANNEL_IDS {

			_, err := mp.s.ChannelMessageSend(
				channel.ChannelID,
				fmt.Sprintf(
					"DEAD (%s): Lab%s -> %s killed by %s",
					discordTimestamp(time.Now(), 0, 0),
					lab,
					boss,
					killer,
				),
			)

			if err != nil {
				log.Printf(
					"Error sending message (labBossHandler Message -kill) to channel (%s). err: %s\n", channel.ChannelID, err,
				)
			}
		}
	}
}

func (mp *MessageParser) cityWarHandler() {
	var (
		owner  string
		city   string
		groups []string
	)

	for _, cwID := range constants.CITY_WAR_CHANNEL_AND_MESSAGE_ID {
		cwMessage, err := mp.s.ChannelMessage(cwID.ChannelID, cwID.MessageID)
		lines := strings.Split(cwMessage.Content, "\n")

		groups = cityWarWinRegex.FindStringSubmatch(mp.m.Content)
		owner = groups[1]
		city = groups[2]

		if err != nil {
			log.Printf(
				"Failed to get message from channel: (%s) message: (%s) (cityWarHandler)- err: %s.\n",
				cwID.ChannelID, cwID.MessageID, err,
			)
			continue
		}

		lines[1] = fmt.Sprintf("**CW Results**\t(%s)", discordTimestamp(time.Now().UTC(), 0, 0))

		switch city {
		case "PC":
			lines[2] = fmt.Sprintf("PC:\t%s", owner)
		case "AC":
			lines[3] = fmt.Sprintf("AC:\t%s", owner)
		case "DC":
			lines[4] = fmt.Sprintf("DC:\t%s", owner)
		case "BI":
			lines[5] = fmt.Sprintf("BI:\t%s", owner)
		}

		lines[7] = fmt.Sprintf("**Next City War**: %s", discordTimestamp(time.Now().UTC(), 4, 10))

		_, err = mp.s.ChannelMessageEdit(cwID.ChannelID, cwID.MessageID, strings.Join(lines, "\n"))

		if err != nil {
			log.Printf(
				"Failed to update message from channel (%s) message (%s) (cityWarHandler)- err: %s.\n",
				cwID.ChannelID, cwID.MessageID, err,
			)
		}
	}
}

// ----- HELPERS -----
func (mp *MessageParser) rolesGeneratorRates(rateType string) ServerRoleMap {

	var (
		roleMap ServerRoleMap = make(map[string][]string)
		allIds  []constants.RoleIDs
		tenXIds []constants.RoleIDs
	)

	if strings.Contains(mp.m.Content, "finished") {
		// if the message has "finished" in it, don't add notifications (roles)
		return roleMap
	}

	stringRate := argParser(getRateValueRegex.FindString(mp.m.Content), "x")
	rate, err := strconv.Atoi(stringRate)

	switch rateType {
	case "drop":
		allIds = constants.DROP_RATE_ROLE_IDS
		tenXIds = constants.DROP_RATE_10X_ROLE_IDS
	case "gold":
		allIds = constants.GOLD_MULTIPLIER_RATE_ROLE_IDS
		tenXIds = constants.GOLD_MULTIPLIER_RATE_10X_ROLE_IDS
	case "dbSoc":
		allIds = constants.DB_SOC_RATE_ROLE_IDS
		tenXIds = constants.DB_SOC_RATE_10X_ROLE_IDS
	case "metSoc":
		allIds = constants.MET_SOC_RATE_ROLE_IDS
		tenXIds = constants.MET_SOC_RATE_10X_ROLE_IDS
	}

	if err != nil {
		log.Printf("Error parsing drop rate value (%s).\n", stringRate)
		return roleMap
	}

	for _, allRateIds := range allIds {
		roleMap[allRateIds.ServerID] = append(roleMap[allRateIds.ServerID], fmt.Sprintf("<@&%s>", allRateIds.RoleID))
	}

	if rate >= 10 {
		for _, tenXRateIds := range tenXIds {
			roleMap[tenXRateIds.ServerID] = append(roleMap[tenXRateIds.ServerID], fmt.Sprintf("<@&%s>", tenXRateIds.RoleID))
		}
	}

	return roleMap
}

func (mp *MessageParser) roleGenerator(event string, role *ServerRoleMap) ServerRoleMap {
	var roleMap ServerRoleMap

	if role == nil {
		roleMap = make(map[string][]string)
	} else {
		roleMap = *role
	}

	switch event {
	case "labBoss":
		for _, role := range constants.LAB_BOSS_ROLE_IDS {
			roleMap[role.ServerID] = append(roleMap[role.ServerID], fmt.Sprintf("<@&%s>", role.RoleID))
		}
	case "gibbon":
		for _, role := range constants.LAB_BOSS_GIBBON_ROLE {
			roleMap[role.ServerID] = append(roleMap[role.ServerID], fmt.Sprintf("<@&%s>", role.RoleID))
		}
	case "nagaLord":
		for _, role := range constants.LAB_BOSS_NAGA_LORD_ROLE {
			roleMap[role.ServerID] = append(roleMap[role.ServerID], fmt.Sprintf("<@&%s>", role.RoleID))
		}
	case "talon":
		for _, role := range constants.LAB_BOSS_TALON_ROLE {
			roleMap[role.ServerID] = append(roleMap[role.ServerID], fmt.Sprintf("<@&%s>", role.RoleID))
		}
	case "howler":
		for _, role := range constants.LAB_BOSS_HOWLER_ROLE {
			roleMap[role.ServerID] = append(roleMap[role.ServerID], fmt.Sprintf("<@&%s>", role.RoleID))
		}
	}

	return roleMap
}

func (mp *MessageParser) isAdminUser() bool {
	return mp.m.Author != nil && slices.Contains(constants.ADMIN_USERS, mp.m.Author.ID)
}

func argParser(s string, separator string) string {
	splitRes := strings.Split(s, separator)
	// trim is for cleanup
	return strings.Trim(splitRes[1], `"`)
}

func addTimeOffset(t time.Time, hours int, minutes int, format string) string {
	// Add hours and minutes with correct units
	newTime := t.
		Add(time.Duration(hours) * time.Hour).
		Add(time.Duration(minutes) * time.Minute)

	// Load the EST/EDT (America/New_York) timezone
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		// fallback: just return UTC formatted
		return newTime.Format(format)
	}

	// Convert to EST/EDT and format
	return newTime.In(loc).Format(format)
}

func discordTimestamp(t time.Time, hours, minutes int) string {
	newTime := t.
		Add(time.Duration(hours) * time.Hour).
		Add(time.Duration(minutes) * time.Minute)

	return fmt.Sprintf("<t:%d:t>", newTime.Unix())
}
