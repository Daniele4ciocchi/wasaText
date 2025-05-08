package api

import (
	"fmt"
	"net/http"
	"strings"
)

type Token struct {
	Token string `json:"Authorization"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Conversation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID      int    `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

// controlla se il token è valido e se si trova all'interno del db
// se non è valido restituisce un errore
// se è valido restituisce l'id dell'utente
func checkAuth(rt *_router, r *http.Request) (int, error) {
	var token Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return -1, fmt.Errorf("missing Authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return -1, fmt.Errorf("invalid Authorization header format")
	}
	token.Token = strings.TrimPrefix(authHeader, "Bearer ")
	if token.Token == "" {
		return -1, fmt.Errorf("empty token")
	}

	id, err := rt.db.CheckToken(token.Token)
	if err != nil {
		return -1, fmt.Errorf("invalid token: %w", err)
	}
	return id, nil
}
