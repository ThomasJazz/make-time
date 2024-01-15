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
	activeGames map[string]GameState
	mtx         sync.Mutex
)

const dumpPath = "/data/[pid].json"

func HandleBlackJack(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := util.ParseLine(s, m)
	if len(args) > 3 {
		s.ChannelMessageSend(m.ChannelID, "Too many of arguments provided")
	}

	game := LoadOrCreateGameState(m.Author.ID)
	var fullResponse strings.Builder
	var actionResult *Action

	if util.Argument(args[1]) == Bet {
		betValue, err := strconv.Atoi(args[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "invalid bet value")
			return
		}
		// Todo validate that bet value is less than player balance
		game.Pot += betValue

		return
	} else {
		action := args[1]
		var err error

		actionResult, err = PlayTurn(action, &game)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error occurred while performing action")
			saveGameAsJson(game)
			return
		}
	}

	// Add to the response as we
	fullResponse.WriteString(formatActionResultString(*actionResult) + "\n")

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
		s.ChannelMessageSend(m.ChannelID, getPlayerTableView(game, false))
	}
}

func PlayTurn(actionStr string, game *GameState) (*Action, error) {
	action := util.Argument(actionStr)
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
	case Stay:
		actionResult.Result = Under
		playerSum := getHandSum(game.PlayerHand, true)

		// Do dealer actions until bust, stay, or blackjack
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

func LoadOrCreateGameState(userId string) GameState {
	var game GameState

	// Either create a new game or load an existing game
	if _, exists := activeGames[userId]; exists {
		newGame, err := loadGameJson(userId)
		if err != nil {
			fmt.Printf("could not load game for user %v", userId)

		}

		copyGame(*newGame, &game)
	} else {
		game = SetupNewGame(userId)
		dealToPlayer(&game)
		dealToPlayer(&game)
		dealToDealer(&game, true)
		dealToDealer(&game, false)
	}

	return game
}

func SetupNewGame(playerId string) GameState {
	suits := []Suit{Spades, Hearts, Diamonds, Clubs}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	var deck []Card

	for _, suit := range suits {
		for i, rank := range ranks {
			deck = append(deck, Card{
				Rank:   rank,
				Suit:   suit,
				Points: util.Min(i+1, 10),
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
	game.DealerHand = append(game.DealerHand, game.Deck[len(game.Deck)-1])
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
