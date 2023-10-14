package main

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	mdl "mastermind/web-service/model"
)

// Non-persistent storage for now
var games = map[uuid.UUID]mdl.Game{}

func main() {
	//- TEST GAME -------------------------
	game := mdl.NewGame()
	test, _ := uuid.Parse("618c3faa-b7ff-4aa9-a618-48a6b88ed0f7")
	games[test] = game
	//- TEST GAME -------------------------

	router := gin.Default()

	// TEST FEATURE
	router.GET("/games", getGames)

	router.POST("/create", newGame)
	router.GET("/games/:token", getGameByToken)
	router.PATCH("/games/:token", updateGameByToken)
	router.DELETE("/games/:token", deleteGameByToken)
	router.Run("localhost:8080")
}

// Create a new game and return a confirmation as response.
func newGame(c *gin.Context) {
	game := mdl.NewGame()
	games[game.Token] = game
	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": "A new game has been created. Good luck!",
		"token":   game.Token.String(),
	})
}

// TEST FEATURE
func getGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, games)
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then return that game
// as a response.
func getGameByToken(c *gin.Context) {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	game, ok := games[token]
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "game not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, game)
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then delete the game
// but don't return any content.
func deleteGameByToken(c *gin.Context) {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	delete(games, token)
	c.IndentedJSON(http.StatusNoContent, nil)
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then update that game
// and return feedback as a response.
func updateGameByToken(c *gin.Context) {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	guess, ok := c.GetPostForm("guess")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	game, ok := games[token]
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "game not found"})
		return
	}
	feedback, err := game.Update(guess)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, feedback)
}
