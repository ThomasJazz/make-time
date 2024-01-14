package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thomasjazz/make-time/util"
)

type ScheduledEvent struct {
	Id        int
	Organizer string
	Attendees []string // Includes the Organizers username and all invitees
	Datetime  time.Time
}

func ScheduleEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	event, err := ParseScheduleCommand(m.Author.Username, m.Content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid command syntax")
		return
	}
	fmt.Println(event)
}

func ParseScheduleCommand(author string, msg string) (*ScheduledEvent, error) {
	var err error

	event := ScheduledEvent{
		Id:        1,
		Organizer: author,
		Attendees: []string{author},
	}

	words := strings.Split(msg, " ")

	for i, word := range words {
		if i == 0 {
			continue
		}

		if strings.HasPrefix(word, util.AtSymbol) {
			event.Attendees = append(event.Attendees, word)
		} else {
			fmt.Println(event)
			err = fmt.Errorf("could not parse word: %s", word)
			return nil, err
		}
	}

	return &event, nil
}
