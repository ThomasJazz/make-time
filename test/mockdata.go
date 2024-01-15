package test

import (
	"github.com/bwmarrin/discordgo"
)

func MockDiscordMessage() (discordgo.Session, discordgo.MessageCreate) {
	session := discordgo.Session{}
	message := discordgo.MessageCreate{
		Message: &discordgo.Message{ChannelID: "channelid", Author: &discordgo.User{ID: "userid"}},
	}

	return session, message
}
