package test

import (
	"testing"
	"time"

	"github.com/thomasjazz/make-time/commands"
)

func TestParseScheduleCommand(t *testing.T) {
	command := "!schedule @Zebrowski @Grimz --date 01-15-2024 --time 17:00"
	author := "@xen0n"
	expected := commands.ScheduledEvent{
		Id:        1,
		Organizer: author,
		Attendees: []string{
			"@xen0n",
			"@Zebrowski",
			"@Grimz",
		},
	}

	eventDate, err := time.Parse("2006-01-02", "2024-01-15")
	if err != nil {
		t.Fatalf("Could not parse date string: 01-15-2024")
	}

	expected.Datetime = eventDate

	result, err := commands.ParseScheduleCommand(author, command)
	if len(result.Attendees) == len(expected.Attendees) {
		t.Fatalf("Number of attendees did not match expected")
	}
}
