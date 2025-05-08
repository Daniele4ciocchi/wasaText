package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	users, err := rt.db.GetUsers()
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

	name := ps.ByName("name")
	id, name, username, err := rt.db.GetUser(name)
	if err != nil {
		http.Error(w, "Errore nel recupero dell'utente", http.StatusInternalServerError)
		return
	}

	user := User{
		ID:       id,
		Name:     name,
		Username: username,
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	id, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	conversations, err := rt.db.GetConversations(id)
	if err != nil {
		http.Error(w, "Errore nel recupero delle conversazioni", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
