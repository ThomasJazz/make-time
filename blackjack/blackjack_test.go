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
		var status Status
		var err error
		status, err = PlayTurn(string(Hit), &game)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(status)
	}

	var status Status
	var err error
	status, err = PlayTurn(string(Stand), &game)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(getPlayerTableView(game, true))
	fmt.Println(status)

}

func TestConvertConst(t *testing.T) {
	arg := util.Argument("asdasd")
	fmt.Println(arg)
}

func TestHandleBlackJackNew(t *testing.T) {
	sess, mess := test.MockDiscordMessage()
	mess.Content = "!blackjack bet 100"

	HandleBlackJack(&sess, &mess)
}

func TestHandleBlackJackLoadedHit(t *testing.T) {
	sess, mess := test.MockDiscordMessage()
	mess.Content = "!blackjack hit"

	HandleBlackJack(&sess, &mess)
}

func TestHandleBlackJackLoadedStand(t *testing.T) {
	sess, mess := test.MockDiscordMessage()
	mess.Content = "!blackjack stand"

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
