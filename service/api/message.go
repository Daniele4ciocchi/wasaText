package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Controlla se il token è valido e se si trova all'interno del db
	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	// Decodifica il messaggio dal corpo della richiesta
	var message utils.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	message.ConversationID, err = strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	//ottengo il mittente dal suo token
	sender, err := rt.db.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	// Controlla se la conversazione esiste nelle conversazioni del mittente
	var conv []utils.Conversation
	conv, err = rt.db.GetConversations(sender.ID)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle conversazioni", http.StatusInternalServerError)
		return
	}

	var found bool
	found = false
	for _, c := range conv {
		print(c.ID)
		if c.ID == message.ConversationID {

			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Conversazione non trovata", http.StatusNotFound)
		return
	}

	//
	message.SenderID = sender.ID

	messageID, err := rt.db.AddMessage(message.SenderID, message.ConversationID, message.Content, message.RepliedMessageID)
	if err != nil {
		http.Error(w, "Errore durante l'invio del messaggio", http.StatusInternalServerError)
		return
	}

	message, err = rt.db.GetMessage(messageID)
	if err != nil {
		http.Error(w, "Errore durante il recupero del messaggio", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Errore durante l'encoding del messaggio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rt *_router) getMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	var conv utils.Conversation
	convID, err := strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	conv, err = rt.db.GetConversation(convID)
	print(conv.ID)
	if err != nil {
		http.Error(w, "Errore durante il recupero della conversazione", http.StatusInternalServerError)
		return
	}

	var messages []utils.Message
	messages, err = rt.db.GetMessages(conv.ID)
	if err != nil {
		http.Error(w, "Errore durante il recupero dei messaggi", http.StatusInternalServerError)
		//printo l'errore
		fmt.Printf("err: %v\n", err)
		return
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Errore durante l'encoding dei messaggi", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
