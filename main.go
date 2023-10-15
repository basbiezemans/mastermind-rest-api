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
	router := newRouter()
	router.Run(":8080")
}

// For testing purposes only.
func addMockGames() {
	tokens := []string{
		"0fd253d0-80dc-42e8-aa0c-b1e9ce84936d",
		"20d245fd-f724-4e1c-a818-04b3dd33ef5d",
	}
	for _, token := range tokens {
		game := mdl.MockGame(uuid.MustParse(token))
		games[game.Token] = game
	}
}

func newRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/create", newGame)
	router.GET("/games/:token", getGameByToken)
	router.PATCH("/games/:token", updateGameByToken)
	router.DELETE("/games/:token", deleteGameByToken)
	return router
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
