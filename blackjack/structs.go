package blackjack

import (
	"github.com/thomasjazz/make-time/util"
)

type Suit string
type Rank string

type Card struct {
	Rank    Rank
	Suit    Suit
	Points  int
	Visible bool
}

type GameState struct {
	PlayerId   string
	Deck       []Card
	DealerHand []Card
	PlayerHand []Card
	Pot        int
}

type Result string
type Status string

type Action struct {
	Action util.Argument
	Result Result
	Status Status
}

// Blackjack possible arguments
var (
	Bet  util.Argument = "bet"
	Hit  util.Argument = "hit"
	Stay util.Argument = "stay"
)

const (
	BetModifier float64 = 1.5
	Hearts      Suit    = "Hearts"
	Spades      Suit    = "Spades"
	Diamonds    Suit    = "Diamonds"
	Clubs       Suit    = "Clubs"

	PlayerBust Result = "Player Bust"
	Under      Result = "Under"
	Blackjack  Result = "Player Blackjack"

	DealerWin  Status = "Dealer Win"
	PlayerWin  Status = "Player win"
	Draw       Status = "Draw"
	InProgress Status = "In-Progress"

	Ace   Rank = "Ace"
	Two   Rank = "Two"
	Three Rank = "Three"
	Four  Rank = "Four"
	Five  Rank = "Five"
	Six   Rank = "Six"
	Seven Rank = "Seven"
	Eight Rank = "Eight"
	Nine  Rank = "Nine"
	Ten   Rank = "Ten"
	Jack  Rank = "Jack"
	Queen Rank = "Queen"
	King  Rank = "King"
)
