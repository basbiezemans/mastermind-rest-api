package internal

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUpdateInvalidGuess(t *testing.T) {
	tests := []struct {
		guess string
		descr string
	}{
		{"", "empty guess"},
		{"123", "guess too short"},
		{"12345", "guess too long"},
		{"0234", "guess contains zero"},
		{"1237", "guess contains seven"},
		{"12a4", "guess contains a letter"},
	}
	game := MockGame(uuid.UUID{})
	for _, test := range tests {
		game.Reset()
		_, err := game.Update(test.guess)
		if err == nil {
			t.Errorf("Update: %s, expected an error, got nil", test.descr)
		}
	}
}

func TestUpdateFeedback(t *testing.T) {
	game := MockGame(uuid.UUID{})
	tests := []struct {
		secret string
		guess  string
		expect string
		descr  string
	}{
		{"1234", "1234", "●●●●", "all digits correct"},
		{"6243", "6225", "●●", "two correct positions"},
		{"5256", "2244", "●", "one correct position"},
		{"1111", "2222", "", "no matching digits"},
		{"6423", "2252", "○", "one present digit"},
		{"6443", "4124", "○○", "two present digits"},
		{"6163", "1136", "●○○", "one correct and two present digits"},
		{"1234", "2134", "●●○○", "two correct and two present digits"},
		{"1234", "2341", "○○○○", "all digits present in wrong positions"},
		{"1234", "1235", "●●●", "three correct positions"},
		{"1234", "1256", "●●", "two correct positions and two absent digits"},
	}
	for _, test := range tests {
		game.Turn = 0
		game.Secret = newCode(test.secret)
		res, err := game.Update(test.guess)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
			continue
		}
		if test.expect != res.Feedback {
			t.Errorf(
				"Feedback: %s, expected %s, got %s",
				test.descr,
				test.expect,
				res.Feedback,
			)
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
