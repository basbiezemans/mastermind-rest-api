package main

import (
	"fmt"
	"mastermind/web-service/model"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func main() {
	model.ConnectDatabase()
	router := newRouter()
	router.Run()
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
	game := model.CreateGame()
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
	game, err := model.GetGame(token)
	if err != nil {
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
	model.DeleteGame(token)
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
	feedback, err := model.UpdateGame(token, guess)
	if err != nil {
		go LogError(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}
	c.IndentedJSON(http.StatusOK, feedback)
}

func LogError(message string) {
	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	fname := "data/error.log"
	logfile, err := os.OpenFile(fname, flags, 0644)
	if err == nil {
		fmt.Fprintln(logfile, message)
	}
}
