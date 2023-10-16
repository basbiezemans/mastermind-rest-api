package model

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUpdateInvalidGuess(t *testing.T) {
	game := MockGame(uuid.UUID{})
	_, err := game.Update("")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

func TestUpdateFeedback(t *testing.T) {
	game := MockGame(uuid.UUID{})
	tests := []struct {
		secret string
		guess  string
		expect string
	}{
		{"1234", "1234", "●●●●"},
		{"6243", "6225", "●●"},
		{"5256", "2244", "●"},
		{"1111", "2222", ""},
		{"6423", "2252", "○"},
		{"6443", "4124", "○○"},
		{"6163", "1136", "●○○"},
		{"1234", "2134", "●●○○"},
	}
	for _, test := range tests {
		game.Secret = Secret{newCode(test.secret)}
		res, err := game.Update(test.guess)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
			continue
		}
		if test.expect != res.Feedback {
			t.Errorf("Feedback, expected %s, got %s", test.expect, res.Feedback)
		}
	}
}

func TestUpdateGameWon(t *testing.T) {
	game := MockGame(uuid.UUID{})
	before := game.Score.CodeBreaker
	res, err := game.Update("1234")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		return
	}
	if !strings.HasPrefix(res.Message, "You won") {
		t.Errorf("Expected code breaker to win but no sigar")
	}
	if game.Turn != 0 {
		t.Errorf("Expected game to be reset after winning")
	}
	after := game.Score.CodeBreaker
	if after != before+1 {
		t.Errorf("Expected code breaker score to be incremented by 1")
	}
}

func TestUpdateGameLost(t *testing.T) {
	game := MockGame(uuid.UUID{})
	game.Turn = 9
	before := game.Score.CodeMaker
	res, err := game.Update("4321")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
		return
	}
	if !strings.HasPrefix(res.Message, "You lost") {
		t.Errorf("Expected code breaker to lose")
	}
	if game.Turn != 0 {
		t.Errorf("Expected game to be reset after game over")
	}
	after := game.Score.CodeMaker
	if after != before+1 {
		t.Errorf("Expected code maker score to be incremented by 1")
	}
}
