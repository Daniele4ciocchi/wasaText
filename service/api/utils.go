package api

import (
	"fmt"
	"net/http"
	"strings"
)

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

func getUserIDFromAuth(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("token mancante")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return "", fmt.Errorf("formato token non valido")
	}

	return strings.TrimPrefix(auth, prefix), nil
}
