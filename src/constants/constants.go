package constants

type CityWarIDs struct {
	ChannelID string
	MessageID string
}

var (
	ASSIGN_ROLE_MESSAGE_IDS = []string{
		"1408951046323175618", // P&O Server
	}

	LISTENING_CHANNEL_IDS = []string{
		"1068927561058488353", // COX event feed channel id
		"1325993204725583967", // P&O general
	}

	RELAY_MESSAGE_CHANNEL_IDS = []string{
		"1408930128205053972", // P&O events feed channel
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
	GOLD_MULTIPLIER_RATE_5X_ROLE_IDS = []string{
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

	CITY_WAR_CHANNEL_AND_MESSAGE_ID = []CityWarIDs{
		{ChannelID: "1409693429482393710", MessageID: "1409701198692352151"}, // P&O Channel/Message
	}
)
