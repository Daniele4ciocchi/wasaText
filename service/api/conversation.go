package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")

	// controllo se il token è valido
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
	var conv utils.Conversation
	var convID int
	var convName string

	// recupero il nome dei due utenti
	user2, err = rt.db.GetUser(name.Name)
	if err != nil {
		http.Error(w, "Utente1 non trovato", http.StatusNotFound)
		return
	}

	user1, err = rt.db.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	if user1.Name == user2.Name {
		http.Error(w, "Non puoi creare una conversazione con te stesso", http.StatusBadRequest)
		return
	}

	// controllo se la conversazione esiste già
	convID, err = rt.db.CheckExistingConversation(user1.ID, user2.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			http.Error(w, "Errore durante il recupero della conversazione", http.StatusInternalServerError)
			return
		}
	}
	if convID != 0 {
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
		return
	}

	// se la conversazione non esiste la creo
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

	for i, conv := range convs {
		if !conv.IsGroup {
			conv.Name = strings.Replace(conv.Name, user.Name, "", -1)
			convs[i] = conv
		}
	}

	if err := json.NewEncoder(w).Encode(convs); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// controllo se il token è valido
	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	// recupero l'utente dal token
	var user utils.User
	user, err = rt.db.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	// recupero l'id della conversazione
	var convID int
	convID, err = strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	conv, err := rt.db.GetConversation(convID)
	if err != nil {
		http.Error(w, "Errore durante il recupero della conversazione", http.StatusInternalServerError)
		return
	}

	if !conv.IsGroup {
		conv.Name = strings.Replace(conv.Name, user.Name, "", -1)
	}

	if err := json.NewEncoder(w).Encode(conv); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
