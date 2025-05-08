package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Errore nella decodifica JSON", http.StatusBadRequest)
		return
	}

	if err := rt.db.AddUser(user.Name, user.Name); err != nil {
		http.Error(w, "Errore nell'aggiunta dell'utente", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
