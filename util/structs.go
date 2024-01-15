package util

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

type Function string
type CommandGroup string

type Command string

const (
	// Static Commands
	CommandHelp   Command = "help"
	CommandGamble Command = "gamble"
	CommandPing   Command = "ping"
	CommandMikey  Command = "mikey"

	CommandScheduleNew Command = "new"

	CommandBalance   Command = "balance"
	CommandBlackJack Command = "blackjack"
	CommandCoinFlip  Command = "flipcoin"

	CommandSchedule Command = "schedule"
)

type Argument string

const (
	// Arguments
	ArgTime Argument = "-t"
)

const (
	MikeyBday = "1995-09-14"
)

type BankBalance struct {
	Id              int
	DiscordMemberId int
	Balance         int64
}
