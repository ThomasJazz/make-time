package lib

import "time"

type ScheduledEvent struct {
	Id        int
	Organizer string
	Attendees []string // Includes the Organizers username and all invitees
	Datetime  time.Time
}

const (
	AtPrefix  = "@"
	CmdPrefix = "!"
)

type CommandGroup string

const (
	CommandSchedule CommandGroup = "schedule"
	CommandGamble   CommandGroup = "gamble"
)

type Command string

const (
	// Static Commands
	CommandHelp Command = "help"

	// Scheduler commands
	CommandScheduleNew    Command = "new"
	CommandScheduleCancel Command = "cancel"

	// Gambling commands
	CommandGambleBet       Command = "bet"
	CommandGambleBlackJack Command = "blackjack"
)

type Argument string

const (
	// Arguments
	ArgTime Argument = "-t"
)
