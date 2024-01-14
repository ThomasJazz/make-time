package gamble

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thomasjazz/make-time/lib"
	"github.com/thomasjazz/make-time/util"
)

func GambleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := util.ParseLine(s, m)[1:]
	fmt.Println(args)
	cmd := lib.Command(args[0])
	fmt.Println(cmd)
	switch cmd {
	case lib.CommandGambleBet:
		return
	case lib.CommandGambleCoinToss:
		result := CoinToss()
		fmt.Printf("Coin toss result: %s", result)
		s.ChannelMessageSend(m.ChannelID, result)
	case lib.CommandGambleBlackJack:
		return
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
