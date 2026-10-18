package main

import (
	"fmt"
	"io"
	"log"
	"mastermind/web-service/internal"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

var errorLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	errorLog, err := getErrorLog()
	if err != nil {
		return fmt.Errorf("create error log: %w", err)
	}
	defer errorLog.Close()

	accessLog, err := getAccessLog()
	if err != nil {
		return fmt.Errorf("create access log: %w", err)
	}
	defer accessLog.Close()

	errorLogger.SetOutput(errorLog)
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	if err := internal.ConnectDatabase(); err != nil {
		logError(err)
		return fmt.Errorf("connect to database: %w", err)
	}

	router := newRouter(accessLog, errorLog)
	if err := router.Run(":8080"); err != nil {
		logError(err)
		return fmt.Errorf("start HTTP server: %w", err)
	}
	return nil
}

func newRouter(accessWriter, errorWriter io.Writer) *gin.Engine {
	router := gin.New()
	router.Use(
		gin.LoggerWithWriter(accessWriter),  // access logs
		gin.RecoveryWithWriter(errorWriter), // panic/recovery logs
	)
	router.POST("/create", newGame)
	router.GET("/", getEndpointInfo)
	router.GET("/games/:token", getGameByToken)
	router.PATCH("/games/:token", updateGameByToken)
	router.DELETE("/games/:token", deleteGameByToken)
	return router
}

// Create a new game and return a confirmation as response.
func newGame(c *gin.Context) {
	game, err := internal.CreateGame()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Game not created",
		})
		logError(err)
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
		logError(err)
		return
	}
	game, err := internal.GetGame(token)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		logError(err)
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
		logError(err)
		return
	}
	numRows := internal.DeleteGame(token)
	if numRows == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Game not found"})
	} else {
		c.IndentedJSON(http.StatusNoContent, nil)
	}
}

// Locate a game whose Token value matches the token
// parameter sent by the client, then update that game
// and return feedback as a response.
func updateGameByToken(c *gin.Context) {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		logError(err)
		return
	}
	guess, ok := c.GetPostForm("guess")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	feedback, err := internal.UpdateGame(token, guess)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		logError(err)
		return
	}
	c.IndentedJSON(http.StatusOK, feedback)
}

func getErrorLog() (*os.File, error) {
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	return os.OpenFile("data/error.log", flags, 0644)
}

func getAccessLog() (*os.File, error) {
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	return os.OpenFile("data/access.log", flags, 0644)
}

func logError(err error) {
	if err == nil {
		return
	}
	// Output(2, ...) identifies the caller of this function
	if writeErr := errorLogger.Output(2, err.Error()); writeErr != nil {
		fmt.Fprintf(os.Stderr, "failed to write error log: %v (original error: %v)\n", writeErr, err)
	}
}
