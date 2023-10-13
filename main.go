package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type result struct {
	Token    uuid.UUID `json:"token"`
	Message  string    `json:"message"`
	Feedback string    `json:"feedback"`
}

type game struct {
	Token            uuid.UUID `json:"token"`
	CreatedOn        time.Time `json:"created_on"`
	CodeMakerScore   uint32    `json:"maker_score"`
	CodeBreakerScore uint32    `json:"breaker_score"`
}

var games map[uuid.UUID]game

func main() {
	// TEST GAME
	token := uuid.New()
	games[token] = game{
		Token:            token,
		CreatedOn:        time.Now(),
		CodeMakerScore:   0,
		CodeBreakerScore: 1,
	}

	router := gin.Default()
	router.POST("/create", newGame)
	router.GET("/games", getGames)
	router.GET("/games/:token", getGameByToken)
	router.PATCH("/games/:token", updateGameByToken)
	router.DELETE("/games/:token", deleteGameByToken)
	router.Run("localhost:8080")
}

// Create a new game and return that game as a response.
func newGame(c *gin.Context) {
	var game = game{
		Token:            uuid.New(),
		CreatedOn:        time.Now(),
		CodeMakerScore:   0,
		CodeBreakerScore: 0,
	}
	games[game.Token] = game
	c.IndentedJSON(http.StatusCreated, game)
}

// Responds with a list of all games.
func getGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, games)
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then return that game
// as a response.
func getGameByToken(c *gin.Context) {
	token := c.Param("token")
	if token, err := uuid.Parse(token); err == nil {
		if game, ok := games[token]; ok {
			c.IndentedJSON(http.StatusOK, game)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "game not found"})
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then update that game
// and return feedback as a response.
func updateGameByToken(c *gin.Context) {
	token := c.Param("token")
	guess := c.Param("guess")
	if token, err := uuid.Parse(token); err == nil {
		if game, ok := games[token]; ok {
			feedback, err := update(game, guess)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
			} else {
				c.IndentedJSON(http.StatusOK, feedback)
			}
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "game not found"})
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then delete the game
// but don't return any content.
func deleteGameByToken(c *gin.Context) {
	token := c.Param("token")
	if token, err := uuid.Parse(token); err == nil {
		if _, ok := games[token]; ok {
			delete(games, token)
			c.IndentedJSON(http.StatusNoContent, nil)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "game not found"})
}

// Update stub
func update(game game, guess string) (result, error) {
	return result{
		Token:    game.Token,
		Message:  "Guess 1 of 10. You guessed: [1,2,3,4]",
		Feedback: "1,0",
	}, nil
}
