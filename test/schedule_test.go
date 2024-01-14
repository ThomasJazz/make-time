package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/thomasjazz/make-time/lib"
	"github.com/thomasjazz/make-time/schedule"
	"github.com/thomasjazz/make-time/util"
)

func TestMikey(t *testing.T) {
	fmt.Println(util.GetMikeyYears())
}
func TestParseScheduleCommand(t *testing.T) {
	command := "!schedule @Zebrowski @Grimz -t \"01-15-2024 17:00\""
	author := "@xen0n"
	expected := lib.ScheduledEvent{
		Id:        1,
		Organizer: author,
		Attendees: []string{
			"@xen0n",
			"@Zebrowski",
			"@Grimz",
		},
	}

	args := []string{command}

	fmt.Println(args)

	eventDate, err := time.Parse("2006-01-02", "2024-01-15")
	if err != nil {
		t.Fatalf("Could not parse date string: 01-15-2024")
	}

	expected.Datetime = eventDate

	result, _ := schedule.ParseScheduleCommand(author, args)
	if len(result.Attendees) == len(expected.Attendees) {
		t.Fatalf("Number of attendees did not match expected")
	}
}

func TestParser(t *testing.T) {
	parsed, _ := util.ParseCommandLine("!gamble flip")
	fmt.Println(parsed)
}
