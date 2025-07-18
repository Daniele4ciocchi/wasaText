package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
)

// tramite la seguente funzione possiamo inviare un messaggio specificandone le informazioni
// all'interno di un json, e l'id della conversazione nel path
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

	// ottengo il mittente dal suo token
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
		if c.ID == message.ConversationID {

			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Conversazione non trovata", http.StatusNotFound)
		return
	}

	message.SenderID = sender.ID

	messageID, err := rt.db.AddMessage(message.SenderID, message.ConversationID, message.Content, message.RepliedMessageID, message.Forwarded)
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

// la funzione serve ad inviare le foto, vengono salvate come messaggi nomrali con un path
// della cartella del backend in cui si trovano
func (rt *_router) sendPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Auth
	id, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Errore nel parsing del form", http.StatusBadRequest)
		return
	}

	// Recupera il file
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Errore nel recupero del file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Recupera l'utente
	user, err := rt.db.GetUserById(id)
	if err != nil {
		http.Error(w, "Errore nel recupero dell'utente", http.StatusInternalServerError)
		return
	}

	// recupero la conversazione
	strconvID := ps.ByName("conversationID")
	convID, err := strconv.Atoi(strconvID)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	var message utils.Message
	message, err = rt.db.GetLastMessage(convID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			message = utils.Message{}
		} else {
			http.Error(w, "Errore durante il recupero del messaggio", http.StatusInternalServerError)
			return
		}
	}

	strmessageID := strconv.Itoa(message.ID + 1)

	path := strmessageID + strconvID

	// Salva il file
	path, err = utils.SaveFile(path, file)
	if err != nil {
		http.Error(w, "Errore nel salvataggio del file", http.StatusInternalServerError)
		return
	}

	// Aggiorna il path nel DB
	_, err = rt.db.AddPhoto(user.ID, convID, path, 0, false)
	if err != nil {
		http.Error(w, "Errore nella modifica della foto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// la funzione ritorna la foto profilo di una conversazione specifica in base a chi è loggato
func (rt *_router) getConversationPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Controlla se il token è valido e se si trova all'interno del db
	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	convID, err := strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	messageID, err := strconv.Atoi(ps.ByName("messageID"))
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	message, err := rt.db.GetMessage(messageID)
	if err != nil {
		http.Error(w, "Errore durante il recupero del messaggio", http.StatusInternalServerError)
		return
	}

	if message.ConversationID != convID {
		http.Error(w, "Messaggio non trovato nella conversazione", http.StatusNotFound)
		return
	}

	if message.Content == "" {
		http.Error(w, "Nessuna foto associata a questo messaggio", http.StatusNotFound)
		return
	}

	file, err := os.Open(message.Content)
	if err != nil {
		http.Error(w, "Foto non trovata", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.WriteHeader(http.StatusOK)

	// Scrive il contenuto del file nella risposta
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Errore durante la scrittura della foto", http.StatusInternalServerError)
		return
	}
}

// tramite questa funzione siamo in grado di eliminare un messaggio
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Controlla se il token è valido e se si trova all'interno del db
	userID, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	var user utils.User
	user, err = rt.db.GetUserById(userID)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	messageID, err := strconv.Atoi(ps.ByName("messageID"))
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	// Controlla se l'utente è il mittente del messaggio
	message, err := rt.db.GetMessage(messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Messaggio non trovato", http.StatusNotFound)
			return
		}
		http.Error(w, "Errore durante il recupero del messaggio", http.StatusInternalServerError)
		return
	}

	if message.SenderID != user.ID {
		http.Error(w, "Non sei il mittente di questo messaggio", http.StatusForbidden)
		return
	}

	err = rt.db.RemoveMessage(messageID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Errore durante la cancellazione del messaggio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// questa funzione ritorna tutti i messaggi di una conversazione
func (rt *_router) getMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	user, err := rt.db.GetUserFromToken(token)
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
	if err != nil {
		http.Error(w, "Errore durante il recupero della conversazione", http.StatusInternalServerError)
		return
	}

	var messages []utils.Message
	messages, err = rt.db.GetMessages(conv.ID)
	if err != nil {
		http.Error(w, "Errore durante il recupero dei messaggi", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Errore durante l'encoding dei messaggi", http.StatusInternalServerError)
		return
	}

	for _, mess := range messages {
		if mess.Status != 2 {
			err = rt.db.SetViewedMessage(user.ID, mess.ID)
			if err != nil {
				http.Error(w, "Errore durante l'aggiornamento del messaggio", http.StatusInternalServerError)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

// questa funzione ritorna l'ultimo messaggio di una conversazione
func (rt *_router) getLastMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	convID, err := strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	var message utils.Message
	message, err = rt.db.GetLastMessage(convID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			message = utils.Message{}
		}
	}

	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Errore durante l'encoding dell'ultimo messaggio", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// la funzione serve a inoltrare un messaggio
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Controlla se il token è valido e se si trova all'interno del db
	id, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	var message utils.Message
	message.ID, err = strconv.Atoi(ps.ByName("messageID"))
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	message, err = rt.db.GetMessage(message.ID)
	if err != nil {
		http.Error(w, "Errore durante il recupero del messaggio", http.StatusInternalServerError)
		return
	}

	type newreciver struct {
		ID int `json:"receiver_id"`
	}

	var reciver newreciver
	if err := json.NewDecoder(r.Body).Decode(&reciver); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	// Controlla se la conversazione esiste nelle conversazioni del mittente
	var conv []utils.Conversation
	conv, err = rt.db.GetConversations(id)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle conversazioni", http.StatusInternalServerError)
		return
	}
	var found bool
	found = false
	for _, c := range conv {
		if c.ID == reciver.ID {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Conversazione non trovata", http.StatusNotFound)
		return
	}

	if message.Photo {
		// Se il messaggio è una foto, aggiungo la foto alla conversazione del ricevente
		_, err = rt.db.AddPhoto(id, reciver.ID, message.Content, message.RepliedMessageID, true)
		if err != nil {
			http.Error(w, "Errore durante l'invio della foto", http.StatusInternalServerError)
			return
		}
		log.Println("Foto inoltrata con successo")
	} else {
		_, err = rt.db.AddMessage(id, reciver.ID, message.Content, message.RepliedMessageID, true)
		if err != nil {
			http.Error(w, "Errore durante l'invio del messaggio", http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Errore durante l'encoding del messaggio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// tramite questa funzione possiamo ottenere nuovi messaggi
func (rt *_router) getNewMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	user, err := rt.db.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	var messages []utils.Message
	messages, err = rt.db.GetNewMessages(user.ID)
	if err != nil {
		http.Error(w, "Errore durante il recupero dei messaggi", http.StatusInternalServerError)
		return
	}

	for _, mess := range messages {
		err = rt.db.SetArrivedMessage(user.ID, mess.ID)
		if err != nil {
			http.Error(w, "Errore durante l'aggiornamento del messaggio", http.StatusInternalServerError)
			log.Println(err)
		}
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Errore durante l'encoding dei messaggi", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
