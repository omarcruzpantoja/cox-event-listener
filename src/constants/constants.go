package constants

type CityWarIDs struct {
	ChannelID string
	MessageID string
}

type LabBosIDs struct {
	ServerID  string
	ChannelID string
	MessageID string
}

type RoleIDs struct {
	ServerID string
	RoleID   string
}

type RelayMessageIDs struct {
	ServerID  string
	ChannelID string
}

var (
	ROGUE_SERVER_ID        = "1094994806511517776"
	BLACK_HOLLOW_SERVER_ID = "1459735738953240712"
	PO_SERVER_ID           = "1325993203836518431"
	// Whenever you define a message to assign roles to users, you'd
	// set the message ID here. So say you have a "React on emoji X for
	// Y role", the id of the message that contains those reactions will be added
	// here. It is important that your bot has permissions to:
	// 1. Manage Roles
	// 2. Listen to reactions
	// 3. Add reactions to messages (if you intend to use the command in this repo)
	// 4, Send messages in channcel (if you intend to use the command in this repo)
	// Otherwise if you intend to use other app for managing roles, then you can skip this
	ASSIGN_ROLE_MESSAGE_IDS = []string{
		"1415116596392890368", // P&O Server Rates
		"1415141373392453702", // P&O Server Lab Boss Spawn
		"1419148638071947405", // P&O Specific Lab Boss Spawn
		"1459981084526776607", // BlackHollow Server Rates
		"1459979919504183461", // BlackHollow Server Lab Boss Spawn
		"1459979979210362901", // BlackHollow Specific Lab Boss Spawn
		"1460844321367523463", // BlackHollow Phoenix Egg Event
	}

	// This is used to identify which channels to listen events. This is for your
	// primary bot.
	LISTENING_CHANNEL_IDS = []string{
		"1068927561058488353", // COX event feed channel id
		"1325993204725583967", // P&O general
	}

	// This is used to identify which channels to listen for events by a "secodary"
	// bot. The secondary bot relays messages to the primary bot. (imagine your primary
	// bot isn't allowed in a server but the secondary is).
	SECOND_BOT_LISTENING_CHANNEL_IDS = []string{
		"1068927561058488353", // COX event feed channel id
	}

	ADMIN_USERS = []string{
		"132992476813197312", // _stacy user
	}

	// Channels to send messages to, this is the channel you create for mentions to be sent.
	RELAY_MESSAGE_CHANNEL_IDS = []RelayMessageIDs{
		{ServerID: PO_SERVER_ID, ChannelID: "1408930128205053972"},           // P&O events feed channel
		{ServerID: ROGUE_SERVER_ID, ChannelID: "1415864173367398441"},        // ROGUE CHANNEL
		{ServerID: BLACK_HOLLOW_SERVER_ID, ChannelID: "1459923502982365488"}, // BLACK HOLLOW CHANNEL
	}

	// Role IDs for `Drop rate all`
	DROP_RATE_ROLE_IDS = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1408913591125540954"},           // P&O server
		{ServerID: ROGUE_SERVER_ID, RoleID: "1387734426502565969"},        // ROGUE server
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977383032524942"}, // BLACK HOLLOW server
	}

	DROP_RATE_10X_ROLE_IDS = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1408913379137032322"},           // P&O server
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977200617914469"}, // BLACK HOLLOW server
	}

	// Role IDs for `Gold multiplier all`
	GOLD_MULTIPLIER_RATE_ROLE_IDS = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1408913674340663296"},           // P&O server
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977452381016196"}, // BLACK HOLLOW server
	}
	GOLD_MULTIPLIER_RATE_10X_ROLE_IDS = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1408913719555129384"},           // P&O server
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977490905956373"}, // BLACK HOLLOW server
	}

	DB_SOC_RATE_ROLE_IDS = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1408914346322690279"},           // P&O server
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977563089670289"}, // BLACK HOLLOW server
	}
	DB_SOC_RATE_10X_ROLE_IDS = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1408914422407364759"},           // P&O server
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459980853978464268"}, // BLACK HOLLOW server
	}

	MET_SOC_RATE_ROLE_IDS     = []RoleIDs{}
	MET_SOC_RATE_10X_ROLE_IDS = []RoleIDs{}

	LAB_BOSS_ROLE_IDS = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1415136091681587201"},           // LAB_BOSS_SPAWN (P&O)
		{ServerID: ROGUE_SERVER_ID, RoleID: "1392926180818292889"},        // ROGUE
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977614394654721"}, // BLACK HOLLOW
	}

	LAB_BOSS_GIBBON_ROLE = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1419146940892975164"},           // LAB_BOSS_SPAWN (P&O)
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977745923575991"}, // LAB_BOSS_SPAWN (BLACK HOLLOW)
	}

	LAB_BOSS_NAGA_LORD_ROLE = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1419147201275494481"},           // LAB_BOSS_SPAWN (P&O)
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977772784025630"}, // LAB_BOSS_SPAWN (BLACK HOLLOW)
	}

	LAB_BOSS_TALON_ROLE = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1419148293853941760"},           // LAB_BOSS_SPAWN (P&O)
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977810918637648"}, // LAB_BOSS_SPAWN (BLACK HOLLOW)
	}

	LAB_BOSS_HOWLER_ROLE = []RoleIDs{
		{ServerID: PO_SERVER_ID, RoleID: "1419147397472194590"},           // LAB_BOSS_SPAWN (P&O)
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1459977838093537310"}, // LAB_BOSS_SPAWN (BLACK HOLLOW)
	}

	CITY_WAR_CHANNEL_AND_MESSAGE_ID = []CityWarIDs{
		{ChannelID: "1409693429482393710", MessageID: "1409701198692352151"}, // P&O Channel/Message
		{ChannelID: "1459981658534187172", MessageID: "1459982183618969894"}, // BLACK HOLLOW Channel/Message
	}

	LAB_BOSS_CHANNEL_AND_MESSAGE_ID = []LabBosIDs{
		{ChannelID: "1409693429482393710", MessageID: "1415121652672630836"}, // P&O Channel/Message
	}

	PHOENIX_EGG_EVENT_ROLE_IDS = []RoleIDs{
		{ServerID: BLACK_HOLLOW_SERVER_ID, RoleID: "1460843028011942019"}, // BLACK HOLLOW server
	}
)
