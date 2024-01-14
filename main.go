package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/thomasjazz/make-time/commands"
	"github.com/thomasjazz/make-time/util"
)

func main() {
	envFile, _ := godotenv.Read(".env")

	token := envFile["DISCORD_TOKEN"]

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	// Declare intent
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Add message handler(s)
	dg.AddHandler(messageCreate)

	// Wait for a CTRL-C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stop

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Printf("Message: %s\n", m.Content)

	// Check the message content and respond accordingly
	switch content := m.Content; {
	case content == "!ping":
		s.ChannelMessageSend(m.ChannelID, "pong!")
	case strings.HasPrefix(content, util.ScheduleCommandPrefix):
		commands.ScheduleEvent(s, m)
	}
}
