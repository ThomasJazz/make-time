package blackjack

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thomasjazz/make-time/util"
)

var (
	activeGames  map[string]GameState
	mtx          sync.Mutex
	sendResponse = make(chan string, 1)
	closeChannel = make(chan struct{})
)

const dumpPath = "/data/[pid].json"

func validateArgs(hasActiveGame bool, args []string) bool {
	switch numArgs := len(args); {
	case numArgs > 3:
		sendResponse <- "Too many of arguments provided"
		return false
	case numArgs == 1:
		sendResponse <- "usage: !blackjack [action]"
		return false
	}

	// This will make sure "!blackhack hit stand" would fail
	actions := 0

	for i := 1; i < len(args); i++ {
		switch val := PlayOption(args[i]); {
		case !isValidPlayOption(string(val)):
			sendResponse <- "Unexpected argument: " + args[i]
			return false
		case val == Bet:
			if i == len(args)-1 {
				sendResponse <- "Invalid syntax. Please use: !blackjack bet [amount]"
				return false
			}
			if _, err := strconv.Atoi(args[i+1]); err != nil {
				sendResponse <- "Bet amount must be integer value"
				return false
			}
			i++ // Increment here so we don't check it on next iteration
		case val == Hit || val == Stand:
			if !hasActiveGame {
				sendResponse <- "No active game found. Please use: !blackjack bet [amount]"
				return false
			}
			if actions > 0 {
				sendResponse <- "Error. Cannot perform multiple actions"
				return false
			}
			actions++
			continue
		}
	}

	return true
}

func HandleBlackJack(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := util.ParseLine(s, m)
	go sendResponseMessage(s, m)

	hasActive := hasActiveGame(m.Author.ID)
	if !validateArgs(hasActive, args) {
		closeChannel <- struct{}{}
		return
	}

	game := LoadOrCreateGameState(m.Author.ID)

	var fullResponse strings.Builder
	var actionResult *Action = &Action{
		Status: InProgress,
	}

	for i := 1; i < len(args); i++ {
		switch PlayOption(args[i]) {
		case Bet:
			// Error is accounted for in arg validation
			betValue, _ := strconv.Atoi(args[i+1])

			// Todo validate that bet value is less than player balance
			game.Pot += betValue

			i++
		case Hit, Stand:
			var err error

			actionResult, err = PlayTurn(args[i], &game)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error occurred while performing action")
				return
			}
		}
	}

	// Add to the response as we

	if actionResult.Status != InProgress {
		pot := endGame(game)
		fullResponse.WriteString(getPlayerTableView(game, true) + "\n")

		if actionResult.Status == PlayerWin {
			fullResponse.WriteString("You won " + strconv.Itoa(pot) + " shmeckles!\n")
			// todo: add winnings to player balance in DB
		} else {
			fullResponse.WriteString("You lost " + strconv.Itoa(pot) + " shmeckles!\n")
		}
	} else {
		fmt.Println(getPlayerTableView(game, false))
		//s.ChannelMessageSend(m.ChannelID, getPlayerTableView(game, false))
	}

	closeChannel <- struct{}{}
}

func PlayTurn(actionStr string, game *GameState) (*Action, error) {
	action := PlayOption(actionStr)
	actionResult := Action{
		Action: action,
	}

	switch action {
	case Hit:
		dealToPlayer(game)
		handSum := getHandSum(game.PlayerHand, true)
		dealerSum := getHandSum(game.DealerHand, true)
		if handSum == 21 {
			actionResult.Result = Blackjack

			if dealerSum != 21 {
				actionResult.Status = Draw
			} else {
				actionResult.Status = PlayerWin
			}
		} else if handSum > 21 {
			actionResult.Result = PlayerBust
			actionResult.Status = DealerWin
		}
	case Stand:
		actionResult.Result = Under
		playerSum := getHandSum(game.PlayerHand, true)

		// Do dealer actions until bust, stand, or blackjack
		for doDealerAction(game) {
		}

		dealerSum := getHandSum(game.DealerHand, true)

		// Determine winner
		if playerSum > dealerSum {
			actionResult.Status = PlayerWin
		} else if playerSum == dealerSum {
			actionResult.Status = Draw
		} else {
			actionResult.Status = DealerWin
		}
	default:
		return &actionResult, errors.New("invalid action")
	}

	return &actionResult, nil
}

func LoadOrCreateGameState(playerId string) GameState {
	var game GameState

	// Either create a new game or load an existing game
	if hasActiveGame(playerId) {
		newGame, err := loadGameJson(playerId)
		if err != nil {
			fmt.Printf("could not load game for user %v", playerId)

		}

		copyGame(*newGame, &game)
	} else {
		game = StartGame(playerId)
	}

	return game
}

func StartGame(playerId string) GameState {
	game := SetupNewGame(playerId)
	dealToPlayer(&game)
	dealToPlayer(&game)
	dealToDealer(&game, true)
	dealToDealer(&game, false)

	return game
}

func SetupNewGame(playerId string) GameState {
	suits := []Suit{Spades, Hearts, Diamonds, Clubs}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	var deck []Card

	for _, suit := range suits {
		for i, rank := range ranks {
			deck = append(deck, Card{
				Rank:    rank,
				Suit:    suit,
				Points:  util.Min(i+1, 10),
				Visible: true,
			})
		}
	}

	shuffleDeck(deck)

	game := GameState{
		PlayerId:   playerId,
		Deck:       deck,
		DealerHand: []Card{},
		PlayerHand: []Card{},
	}

	return game
}

func shuffleDeck(deck []Card) {
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}

func doDealerAction(game *GameState) bool {
	handSum := getHandSum(game.DealerHand, true)

	if handSum < 17 {
		dealToDealer(game, true)
		return true
	}

	return false
}

func getHandSum(cards []Card, countHidden bool) int {
	sumAceEleven := 0
	sumAceOne := 0
	foundAce := false

	for _, card := range cards {
		if !countHidden && !card.Visible {
			continue
		}

		if card.Rank == Ace {
			sumAceEleven += 11
			sumAceOne += 1
			foundAce = true
		} else {
			sumAceEleven += card.Points
			sumAceOne += card.Points
		}
	}

	if foundAce && sumAceEleven > 21 {
		return sumAceOne
	}

	return sumAceEleven
}

// Deal card to player and remove the card from the deck
func dealToPlayer(game *GameState) {
	game.PlayerHand = append(game.PlayerHand, game.Deck[len(game.Deck)-1])
	game.Deck = game.Deck[:len(game.Deck)-1]
}

func dealToDealer(game *GameState, visible bool) {
	// Use the visibility provided
	card := game.Deck[len(game.Deck)-1]
	card.Visible = visible

	game.DealerHand = append(game.DealerHand, card)
	game.Deck = game.Deck[:len(game.Deck)-1]
}

func saveGameAsJson(game GameState) error {
	jsonData, err := json.Marshal(game)
	if err != nil {
		return err
	}
	filepath := getFilePath(game.PlayerId)

	err = os.WriteFile(filepath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadGameJson(playerId string) (*GameState, error) {
	jsonData, err := os.ReadFile(getFilePath(playerId))
	if err != nil {
		return nil, err
	}

	var game GameState
	err = json.Unmarshal(jsonData, &game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func getFilePath(playerId string) string {
	return strings.Replace(dumpPath, "[pid]", playerId, 1)
}

func endGame(game GameState) int {
	os.Remove(getFilePath(game.PlayerId))
	return int(float64(game.Pot) * BetModifier)
}

func hasActiveGame(playerId string) bool {
	if _, err := os.Stat(getFilePath(playerId)); err == nil {
		return true
	} else {
		return false
	}
}

func copyGame(copy GameState, copyTo *GameState) {
	copyTo.Deck = copy.Deck
	copyTo.DealerHand = copy.DealerHand
	copyTo.PlayerHand = copy.PlayerHand
	copyTo.PlayerId = copy.PlayerId
	copyTo.Pot = copy.Pot
}

func formatActionResultString(actionResult Action) string {
	return "ActionResult ~~ result: " + string(actionResult.Result) + ", status: " + string(actionResult.Status)

}

func getPlayerTableView(game GameState, showHidden bool) string {
	var board strings.Builder
	dealerSum := getHandSum(game.DealerHand, showHidden)
	playerSum := getHandSum(game.PlayerHand, true)

	board.WriteString("DEALER showing: ")
	for i, card := range game.DealerHand {
		if showHidden {
			board.WriteString(string(card.Rank) + " of " + string(card.Suit))
			if i < len(game.DealerHand)-1 {
				board.WriteString(", ")
			}
		} else if card.Visible {
			board.WriteString(string(card.Rank) + " of " + string(card.Suit) + " ")
		}
	}
	board.WriteString("\n\tSum: " + strconv.Itoa(dealerSum))

	board.WriteString("\nPLAYER showing: ")
	for i, card := range game.PlayerHand {
		board.WriteString(string(card.Rank) + " of " + string(card.Suit))
		if i < len(game.PlayerHand)-1 {
			board.WriteString(", ")
		}
	}

	board.WriteString("\n\tSum: " + strconv.Itoa(playerSum))

	return board.String()
}

func isValidPlayOption(option string) bool {
	switch PlayOption(option) {
	case Bet, Hit, Stand:
		return true
	}
	return false
}

func sendResponseMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	for loop := true; loop; {
		select {
		case val := <-sendResponse:
			s.ChannelMessageSend(m.ChannelID, val)
		case <-closeChannel:
			loop = false
		}
	}
}
