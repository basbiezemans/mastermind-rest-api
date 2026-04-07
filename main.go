package main

import (
	"fmt"
	"io"
	"mastermind/web-service/model"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func main() {
	initErrorWriter()
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	model.ConnectDatabase()
	router := newRouter()
	router.Run()
}

func newRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/create", newGame)
	router.GET("/", getEndpointInfo)
	router.GET("/games/:token", getGameByToken)
	router.PATCH("/games/:token", updateGameByToken)
	router.DELETE("/games/:token", deleteGameByToken)
	return router
}

// Create a new game and return a confirmation as response.
func newGame(c *gin.Context) {
	game, err := model.CreateGame()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Game not created",
		})
		go logError(err.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": "A new game has been created. Good luck!",
		"token":   game.Token.String(),
	})
}

func getEndpointInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "This is the Mastermind code-breaking game.",
		"version": "v0.1-" + gin.Mode(),
	})
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then return that game
// as a response.
func getGameByToken(c *gin.Context) {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		go logError(err.Error())
		return
	}
	game, err := model.GetGame(token)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		go logError(err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, game.Info())
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then delete the game
// but don't return any content.
func deleteGameByToken(c *gin.Context) {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		go logError(err.Error())
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		go logError(err.Error())
		return
	}
	guess, ok := c.GetPostForm("guess")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	feedback, err := model.UpdateGame(token, guess)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		go logError(err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, feedback)
}

func initErrorWriter() {
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	file, err := os.OpenFile("data/error.log", flags, 0644)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	defer file.Close()
	gin.DefaultErrorWriter = io.MultiWriter(file)
}

func logError(message string) {
	t := time.Now()
	message = fmt.Sprintf("%s | %s\n", t.Format(time.RFC3339), message)
	_, err := gin.DefaultErrorWriter.Write([]byte(message))
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}
