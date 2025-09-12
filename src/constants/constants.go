package constants

type CityWarIDs struct {
	ChannelID string
	MessageID string
}

type LabBosIDs struct {
	ChannelID string
	MessageID string
}

var (
	ASSIGN_ROLE_MESSAGE_IDS = []string{
		"1415116596392890368", // P&O Server Rates
		"1415141373392453702", // P&O Server Lab Boss Spawn
	}

	LISTENING_CHANNEL_IDS = []string{
		"1068927561058488353", // COX event feed channel id
		"1325993204725583967", // P&O general
	}

	ACCOUNT_LISTENING_CHANNEL_IDS = []string{
		"1068927561058488353", // COX event feed channel id
	}

	RELAY_MESSAGE_CHANNEL_IDS = []string{
		"1408930128205053972", // P&O events feed channel
		"1415864173367398441", // ROGUE CHANNEL
	}

	DROP_RATE_ROLE_IDS = []string{
		"1408913591125540954", // P&O server
	}
	DROP_RATE_10X_ROLE_IDS = []string{
		"1408913379137032322", // P&O server
	}

	GOLD_MULTIPLIER_RATE_ROLE_IDS = []string{
		"1408913674340663296", // P&O server
	}
	GOLD_MULTIPLIER_RATE_10X_ROLE_IDS = []string{
		"1408913719555129384", // P&O server
	}

	DB_SOC_RATE_ROLE_IDS = []string{
		"1408914346322690279", // P&O server
	}
	DB_SOC_RATE_10X_ROLE_IDS = []string{
		"1408914422407364759", // P&O server
	}

	MET_SOC_RATE_ROLE_IDS     = []string{}
	MET_SOC_RATE_10X_ROLE_IDS = []string{}

	LAB_BOSS_ROLE_IDS = []string{
		"1415136091681587201", // LAB_BOSS_SPAWN
		"1392926180818292889",
	}

	CITY_WAR_CHANNEL_AND_MESSAGE_ID = []CityWarIDs{
		{ChannelID: "1409693429482393710", MessageID: "1409701198692352151"}, // P&O Channel/Message
	}

	LAB_BOSS_CHANNEL_AND_MESSAGE_ID = []LabBosIDs{
		{ChannelID: "1409693429482393710", MessageID: "1415121652672630836"}, // P&O Channel/Message
	}
)
