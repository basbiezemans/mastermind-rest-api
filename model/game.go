package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	. "github.com/basbiezemans/gofunctools/functools"
)

type Game struct {
	CreatedOn time.Time `json:"created_on"`
	Token     uuid.UUID `json:"token"`
	Score     Score     `json:"score"`
	Turn      uint8     `json:"game_turn"`
	Secret    Secret    `json:"secret_code"`
}

type Score struct {
	CodeMaker   uint32 `json:"code_maker"`
	CodeBreaker uint32 `json:"code_breaker"`
}

type Result struct {
	Token    uuid.UUID `json:"token"`
	Message  string    `json:"message"`
	Guess    string    `json:"guess"`
	Feedback string    `json:"feedback"`
}

type Feedback struct {
	Correct int
	Present int
}

type NumPresent struct {
	Tally  int
	Digits []rune
}

type GameInfo struct {
	CreatedOn time.Time `json:"created_on"`
	Token     uuid.UUID `json:"token"`
	Score     Score     `json:"score"`
}

func NewGame() Game {
	return Game{
		Token:     uuid.New(),
		CreatedOn: time.Now(),
		Turn:      0,
		Score:     NewScore(),
		Secret:    NewSecret(),
	}
}

func MockGame(token uuid.UUID) Game {
	return Game{
		Token:     token,
		CreatedOn: time.Now(),
		Turn:      0,
		Score:     NewScore(),
		Secret:    Secret{Code: newCode("1234")},
	}
}

func NewScore() Score {
	return Score{
		CodeMaker:   0,
		CodeBreaker: 0,
	}
}

// Takes a secret code and a guess, and returns feedback that
// shows how many digits are correct and/or present in the guess.
func NewFeedback(secret Code, guess Code) Feedback {
	var pairs = ZipWith(NewPair, secret.Digits, guess.Digits)
	return Feedback{
		Correct: numCorrect(pairs),
		Present: numPresent(Unequal(pairs)),
	}
}

func numCorrect(pairs []Pair[rune]) int {
	return len(pairs) - len(Unequal(pairs))
}

func count(np NumPresent, r rune) NumPresent {
	if digits, ok := Remove(r, np.Digits); ok {
		np.Tally += 1
		np.Digits = digits
		return np
	}
	return np
}

func numPresent(pairs []Pair[rune]) int {
	secret, guess := UnzipWith(Unpair, pairs)
	return FoldLeft(count, NumPresent{0, secret}, guess).Tally
}

func (f Feedback) String() string {
	return strings.Repeat("●", f.Correct) + strings.Repeat("○", f.Present)
}

// Update the game and return a result
func (g *Game) Update(guess string) (Result, error) {
	code, err := CodeFromString(guess)
	if err != nil {
		return Result{}, err
	}
	feedback := NewFeedback(g.Secret.Code, code)
	// Increment turn count
	g.Turn += 1
	// Award a point if applicable
	isCorrectGuess := feedback.Correct == 4
	isGameOver := g.Turn == 10
	message := fmt.Sprintf("Guess %d of 10. You guessed: %s", g.Turn, guess)
	format := "%s The current score is %d (You) vs %d (CodeMaker)."
	if isCorrectGuess {
		g.Score.CodeBreaker += 1
		points := g.Score
		message = fmt.Sprintf(format, "You won!", points.CodeBreaker, points.CodeMaker)
		g.Reset()
	} else if isGameOver {
		g.Score.CodeMaker += 1
		points := g.Score
		message = fmt.Sprintf(format, "You lost.", points.CodeBreaker, points.CodeMaker)
		g.Reset()
	}
	return Result{
		Message:  message,
		Token:    g.Token,
		Guess:    guess,
		Feedback: feedback.String(),
	}, nil
}

// Reset the game
func (g *Game) Reset() {
	g.Turn = 0
	g.Secret = NewSecret()
}

func (g *Game) Info() GameInfo {
	return GameInfo{
		CreatedOn: g.CreatedOn,
		Token:     g.Token,
		Score:     g.Score,
	}
}
