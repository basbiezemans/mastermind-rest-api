package model

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	Token     uuid.UUID `json:"token"`
	CreatedOn time.Time `json:"created_on"`
	Turn      uint8     `json:"-"`
	Score     Score     `json:"score"`
	Secret    Secret    `json:"-"`
}

type Score struct {
	CodeMaker   uint32 `json:"code_maker"`
	CodeBreaker uint32 `json:"code_breaker"`
}

type Result struct {
	Token    uuid.UUID `json:"token"`
	Message  string    `json:"message"`
	Feedback []int     `json:"feedback"`
}

func NewGame() Game {
	return Game{
		Token:     uuid.New(),
		CreatedOn: time.Now(),
		Turn:      1,
		Score:     Score{0, 0},
		Secret:    NewSecret(),
	}
}

func MockGame(token uuid.UUID) Game {
	return Game{
		Token:     token,
		CreatedOn: time.Now(),
		Turn:      1,
		Score:     Score{0, 0},
		Secret:    NewSecret(),
	}
}

func NewScore() Score {
	return Score{
		CodeMaker:   0,
		CodeBreaker: 0,
	}
}

// Stub: update the game and return a result
func (g Game) Update(guess string) (Result, error) {
	code, err := CodeFromString(guess)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Message:  "Guess 1 of 10. You guessed: " + code.String(),
		Token:    g.Token,
		Feedback: []int{1, 0},
	}, nil
}
