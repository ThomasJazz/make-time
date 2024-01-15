package gamble

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thomasjazz/make-time/util"
)

func GambleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := util.ParseLine(s, m)

	cmd := util.Command(args[0])
	fmt.Printf("Gamble command received: %s\n", cmd)

	switch cmd {
	case util.CommandCoinFlip:
		result := CoinToss()
		fmt.Printf("Coin toss result: %s", result)
		s.ChannelMessageSend(m.ChannelID, result)
	}
}

func CoinToss() string {
	result := util.GetRand(2)

	if result == 0 {
		return "Heads"
	} else {
		return "Tails"
	}
}
