package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if _, err := getUserIDFromAuth(r); err != nil {
		http.Error(w, "Token mancante", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	usernames, err := rt.db.GetAllUsernames()
	if err != nil {
		http.Error(w, "Errore nel recupero dei nomi utente", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(usernames); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) addUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var user struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Errore nella decodifica JSON", http.StatusBadRequest)
		return
	}

	if err := rt.db.AddUser(user.Name, user.Username); err != nil {
		http.Error(w, "Errore nell'aggiunta dell'utente", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rt *_router) SetMyUsername(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	name, err := getUserIDFromAuth(r)
	if err != nil {
		http.Error(w, "Token mancante", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Errore nella decodifica JSON", http.StatusBadRequest)
		return
	}

	if err := rt.db.SetMyUsername(name, payload.Username); err != nil {
		http.Error(w, "Errore nell'aggiornamento del nome utente", http.StatusInternalServerError)
		return
	}
	// Restituisci una risposta con name e username aggiornati
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"name": payload.Username})
}
