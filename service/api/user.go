package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
)

type userResponse struct {
	Name string `json:"name"`
}

// controlla se il token è valido e se si trova all'interno del db
// se non è valido restituisce un errore
// se è valido restituisce l'id dell'utente
func checkAuth(r *_router, req *http.Request) (int, error) {
	var token utils.Token
	authHeader := req.Header.Get("Authorization")
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

	id, err := r.db.CheckToken(token.Token)
	if err != nil {
		return -1, fmt.Errorf("invalid token: %w", err)
	}
	return id, nil
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	var users []utils.User

	users, err = rt.db.GetUsers()
	if err != nil {
		http.Error(w, "Errore nel recupero degli utenti", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var user utils.User

	name := ps.ByName("userID")
	user, err := rt.db.GetUser(name)
	if err != nil {
		http.Error(w, "Errore nel recupero dell'utente", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
}
