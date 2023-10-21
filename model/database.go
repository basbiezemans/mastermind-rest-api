package model

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	. "github.com/basbiezemans/gofunctools/functools"
)

type GameState struct {
	Token   uuid.UUID `gorm:"type:uuid;primaryKey"`
	BinData json.RawMessage
}

var db *gorm.DB

func ConnectDatabase() {
	db = connect("data/mastermind.db")
	createGameStateIfNotExists([]Game{})
}

func ConnectMockDatabase() {
	dirpath := "data"
	_, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		dirpath = "../data"
	}
	db = connect(dirpath + "/mock.db")
	createGameStateIfNotExists(getMockGames())
}

func connect(dsn string) *gorm.DB {
	var conn, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}
	return conn
}

// Create a game_states table in case it doesn't exist.
// Pre-populate the table if needed.
func createGameStateIfNotExists(games []Game) {
	if !db.Migrator().HasTable(GameState{}) {
		db.AutoMigrate(GameState{})
	}
	if len(games) > 0 {
		db.Exec("DELETE FROM game_states")
		db.Create(Map(Game.Convert, games))
	}
}

func getMockGames() []Game {
	tokens := []string{
		"0fd253d0-80dc-42e8-aa0c-b1e9ce84936d",
		"20d245fd-f724-4e1c-a818-04b3dd33ef5d",
	}
	return Map(Compose(MockGame, uuid.MustParse), tokens)
}

func (g Game) Convert() GameState {
	bytes, err := json.Marshal(g)
	if err != nil {
		return GameState{}
	}
	return GameState{g.Token, bytes}
}

func (g GameState) Convert() Game {
	bytes := []byte(g.BinData)
	var game Game
	err := json.Unmarshal(bytes, &game)
	if err != nil {
		return Game{}
	}
	return game
}

func CreateGame() Game {
	game := NewGame()
	db.Create(game.Convert())
	return game
}

func GetGame(token uuid.UUID) (Game, error) {
	var state GameState
	var result = db.Find(&state, token)
	if result.Error != nil {
		return Game{}, result.Error
	}
	if result.RowsAffected == 0 {
		return Game{}, errors.New("game not found")
	}
	return state.Convert(), nil
}

func UpdateGame(token uuid.UUID, guess string) (Result, error) {
	game, err := GetGame(token)
	if err != nil {
		return Result{}, err
	}
	result, err := game.Update(guess)
	if err != nil {
		return Result{}, err
	}
	response := db.Save(game.Convert())
	if response.Error != nil {
		return Result{}, response.Error
	}
	return result, nil
}

func DeleteGame(token uuid.UUID) bool {
	return db.Delete(GameState{}, token).RowsAffected > 0
}
