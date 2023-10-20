package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestMain(t *testing.T) {
	ConnectMockDatabase()
}

func TestCreateGame(t *testing.T) {
	token := CreateGame().Token
	_, err := GetGame(token)
	if err != nil {
		t.Error(err)
	}
}

func TestGetGameSuccess(t *testing.T) {
	token := "0fd253d0-80dc-42e8-aa0c-b1e9ce84936d"
	_, err := GetGame(uuid.MustParse(token))
	if err != nil {
		t.Error(err)
	}
}

func TestGetGameFailure(t *testing.T) {
	token := "11111111-2222-3333-4444-555555555555"
	_, err := GetGame(uuid.MustParse(token))
	if err == nil {
		t.Error("expected game not found")
	}
}

func TestUpdateGameSuccess(t *testing.T) {
	token := "0fd253d0-80dc-42e8-aa0c-b1e9ce84936d"
	game, err := GetGame(uuid.MustParse(token))
	if err != nil {
		t.Error(err)
	}
	_, err = game.Update("1234")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateGameFailure(t *testing.T) {
	token := "0fd253d0-80dc-42e8-aa0c-b1e9ce84936d"
	game, err := GetGame(uuid.MustParse(token))
	if err != nil {
		t.Error(err)
	}
	_, err = game.Update("")
	if err == nil {
		t.Error("expected invalid guess error")
	}
}

func TestDeleteGameSuccess(t *testing.T) {
	token := "20d245fd-f724-4e1c-a818-04b3dd33ef5d"
	isDeleted := DeleteGame(uuid.MustParse(token))
	if !isDeleted {
		t.Error("expected game to be deleted")
	}
}

func TestDeleteGameFailure(t *testing.T) {
	token := "11111111-2222-3333-4444-555555555555"
	isDeleted := DeleteGame(uuid.MustParse(token))
	if isDeleted {
		t.Error("expected game NOT to be deleted")
	}
}
