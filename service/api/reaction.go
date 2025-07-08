package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
)

// la funzione serve ad inserire una reazione ad un messaggio
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// auth control
	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	// recupero l'utente dal token
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	var user utils.User
	user, err = rt.db.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Errore nel recupero dell'utente", http.StatusInternalServerError)
		return
	}

	type newReaction struct {
		Content string `json:"content"`
	}

	var reaction newReaction

	messageID, err := strconv.Atoi(ps.ByName("messageID"))
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.AddReaction(user.ID, messageID, reaction.Content)
	if err != nil {
		http.Error(w, "Errore nell'aggiunta della reazione", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// la funzione serve a togliere una reazione ad un messaggio
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// auth control
	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	reactionID, err := strconv.Atoi(ps.ByName("reactionID"))
	if err != nil {
		http.Error(w, "ID reazione non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.RemoveReaction(reactionID)
	if err != nil {
		http.Error(w, "Errore nella rimozione della reazione", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// la funzione ritorna tutte le reazioni di un messaggio
func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// auth control
	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.Atoi(ps.ByName("messageID"))
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	reactions, err := rt.db.GetReactions(messageID)
	if err != nil {
		http.Error(w, "Errore nel recupero delle reazioni", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(reactions); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
}
