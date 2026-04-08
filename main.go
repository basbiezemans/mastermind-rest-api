package main

import (
	"fmt"
	"io"
	"mastermind/web-service/model"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func main() {
	log, err := getErrorLog()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	defer log.Close()
	gin.DefaultErrorWriter = io.Writer(log)
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	err = model.ConnectDatabase()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	router := newRouter()
	router.Use(gin.Recovery())
	router.Run(":8080")
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
		go logError(err)
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
		go logError(err)
		return
	}
	game, err := model.GetGame(token)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		go logError(err)
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
		go logError(err)
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
		go logError(err)
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
		go logError(err)
		return
	}
	c.IndentedJSON(http.StatusOK, feedback)
}

func getErrorLog() (*os.File, error) {
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	file, err := os.OpenFile("data/error.log", flags, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func logError(err error) {
	t := time.Now()
	// skip=0 to get the caller of this function
	_, filename, line, _ := runtime.Caller(0)
	message := fmt.Sprintf("%s | %s:%d | %s\n",
		t.Format(time.RFC3339),
		filename,
		line,
		err.Error(),
	)
	_, err = gin.DefaultErrorWriter.Write([]byte(message))
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}
