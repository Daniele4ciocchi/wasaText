package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type loginRequest struct {
	Name string `json:"name"`
}

type loginResponse struct {
	ID string `json:"id"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	if len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, "Il nome utente deve avere tra 3 e 16 caratteri", http.StatusBadRequest)
		return
	}

	// Check se utente esiste
	_, err := rt.db.GetNameByUsername(req.Name)
	if err != nil {

		if err := rt.db.AddUser(req.Name, req.Name); err != nil {
			http.Error(w, "Errore nella creazione dell'utente", http.StatusInternalServerError)
			return
		}
	}

	// Restituisci "token"/ID
	resp := loginResponse{ID: req.Name}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
