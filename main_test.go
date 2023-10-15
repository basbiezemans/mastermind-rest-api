package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	addMockGames()
}

func TestCreateGame(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create", nil)
	router := newRouter()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	lookup := map[string]string{}
	json.Unmarshal(rec.Body.Bytes(), &lookup)
	if _, ok := lookup["token"]; !ok {
		t.Errorf("missing JSON field: token")
	}
	result, ok := lookup["message"]
	if !ok {
		t.Errorf("missing JSON field: message")
	}
	expect := "A new game has been created. Good luck!"
	assert.Equal(t, expect, result)
}

func TestGetGame(t *testing.T) {
	token := "0fd253d0-80dc-42e8-aa0c-b1e9ce84936d"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/"+token, nil)
	router := newRouter()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetGameNotFound(t *testing.T) {
	token := "11111111-2222-3333-4444-555555555555"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/games/"+token, nil)
	router := newRouter()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateGame(t *testing.T) {
	token := "0fd253d0-80dc-42e8-aa0c-b1e9ce84936d"
	rec := httptest.NewRecorder()
	data := url.Values{}
	data.Set("guess", "1234")
	payload := strings.NewReader(data.Encode())
	req, _ := http.NewRequest("PATCH", "/games/"+token, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router := newRouter()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateGameBadRequest(t *testing.T) {
	token := "0fd253d0-80dc-42e8-aa0c-b1e9ce84936d"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/games/"+token, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router := newRouter()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteGame(t *testing.T) {
	token := "20d245fd-f724-4e1c-a818-04b3dd33ef5d"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/games/"+token, nil)
	router := newRouter()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteGameBadRequest(t *testing.T) {
	token := "invalid-uuid-token"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/games/"+token, nil)
	router := newRouter()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
