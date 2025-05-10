package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	var name utils.Name

	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	var user1 utils.User
	var user2 utils.User

	user1, err = rt.db.GetUser(name.Name)
	if err != nil {
		http.Error(w, "Utente1 non trovato", http.StatusNotFound)
		return
	}

	user2, err = rt.db.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	var convName string
	var convID int
	convName = user1.Name + user2.Name
	convID, err = rt.db.AddConversation(convName, false)
	if err != nil {
		http.Error(w, "Errore durante la creazione della conversazione", http.StatusInternalServerError)
		return
	}

	err = rt.db.AddUserConversation(user1.ID, convID)
	if err != nil {
		http.Error(w, "Errore durante l'aggiunta della conversazione", http.StatusInternalServerError)
		return
	}

	err = rt.db.AddUserConversation(user2.ID, convID)
	if err != nil {
		http.Error(w, "Errore durante l'aggiunta della conversazione", http.StatusInternalServerError)
		return
	}

	var conv utils.Conversation
	conv, err = rt.db.GetConversation(convID)
	if err != nil {
		http.Error(w, "Errore durante il recupero della conversazione", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(conv); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	var user utils.User
	user, err = rt.db.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	var convs []utils.Conversation
	convs, err = rt.db.GetConversations(user.ID)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle conversazioni", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(convs); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
