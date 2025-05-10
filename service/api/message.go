package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	var message utils.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	var user utils.User
	user, err = rt.db.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	message.SenderID = user.ID

	messageID, err := rt.db.AddMessage(message.SenderID, message.ConversationID, message.Content)
	if err != nil {
		http.Error(w, "Errore durante l'invio del messaggio", http.StatusInternalServerError)
		return
	}

	message.ID = messageID

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
