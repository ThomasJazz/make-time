package schedule

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thomasjazz/make-time/lib"
	"github.com/thomasjazz/make-time/util"
)

func ScheduleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := util.ParseLine(s, m)

	event, err := ParseScheduleCommand(m.Author.Username, args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid command syntax")
		return
	}
	fmt.Println(event)
}

func ParseScheduleCommand(author string, args []string) (*lib.ScheduledEvent, error) {
	//var err error

	event := lib.ScheduledEvent{
		Id:        1,
		Organizer: author,
		Attendees: []string{author},
	}

	fmt.Println(args)

	return &event, nil
}
