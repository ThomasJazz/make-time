package blackjack

import (
	"fmt"
	"testing"

	"github.com/thomasjazz/make-time/test"
	"github.com/thomasjazz/make-time/util"
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
	action, err = PlayTurn(string(Stand), &game)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(getPlayerTableView(game, true))
	fmt.Println(action)

}

func TestConvertConst(t *testing.T) {
	arg := util.Argument("asdasd")
	fmt.Println(arg)
}

func TestHandleBlackJack(t *testing.T) {
	sess, mess := test.MockDiscordMessage()
	mess.Content = "!blackjack bet 100"

	HandleBlackJack(&sess, &mess)
}

func TestArgValidator(t *testing.T) {
	input1 := "!blackjack bet 100 hit"
	args, err := util.ParseCommandLine(input1)
	expected := false

	if err != nil {
		t.Fail()
	}

	actual := validateArgs(false, args)

	if expected != actual {
		t.Fail()
	}
}
