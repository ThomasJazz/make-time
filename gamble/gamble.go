package gamble

import (
	"github.com/bwmarrin/discordgo"
	"github.com/thomasjazz/make-time/lib"
	"github.com/thomasjazz/make-time/util"
)

func GambleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := util.ParseLine(s, m)[1:]

	cmd := lib.Command(args[1])
	switch cmd {
	case lib.CommandGambleBet:
		return
	case lib.CommandGambleCoinToss:
		s.ChannelMessageSend(m.ChannelID, CoinToss())
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
