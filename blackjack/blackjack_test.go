package blackjack

import (
	"fmt"
	"testing"
)

func TestGenerateDeck(t *testing.T) {
	game := LoadOrCreateGameState("xenon")
	fmt.Println(getPlayerTableView(game, false))

	for getHandSum(game.PlayerHand, true) < 17 {
		var action *Action
		var err error
		action, err = PlayTurn(string(Hit), &game)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(action)
	}

	var action *Action
	var err error
	action, err = PlayTurn(string(Stay), &game)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(getPlayerTableView(game, true))
	fmt.Println(action)

}
