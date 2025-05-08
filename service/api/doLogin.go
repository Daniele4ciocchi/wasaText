package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// questa funzione permette il login tramite il nome, se è già presente nel db
// allora viene ritornato l'ID dell'utente esistente, altrimenti viene creato
// il nuovo utente e viene ritornato l'ID
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user User

	type Name struct {
		Name string `json:"name"`
	}
	var name Name

	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	if len(name.Name) < 3 || len(name.Name) > 16 {
		http.Error(w, "Il nome utente deve avere tra 3 e 16 caratteri", http.StatusBadRequest)
		return
	}

	var err error
	user.ID, user.Name, user.Username, err = rt.db.GetUser(name.Name)

	if err != nil {
		//crea un nuovo utente
		if err := rt.db.AddUser(name.Name, name.Name); err != nil {
			http.Error(w, "Errore nell'aggiunta dell'utente", http.StatusInternalServerError)
			return
		}

		user.ID, user.Name, user.Username, err = rt.db.GetUser(name.Name)
		if err != nil {
			http.Error(w, "Errore nella ricerca dell'utente", http.StatusInternalServerError)
			return
		}

		err = rt.db.SetToken(user.ID, user.Name)
		if err != nil {
			http.Error(w, "Errore nella creazione del token", http.StatusInternalServerError)
			return
		}

	}
	user.ID, user.Name, user.Username, err = rt.db.GetUser(name.Name)
	if err != nil {
		http.Error(w, "Errore nella ricerca dell'utente", http.StatusInternalServerError)
		return
	}

	var token Token
	token.Token, err = rt.db.GetToken(user.ID)
	if err != nil {
		http.Error(w, "Errore nella ricerca del token", http.StatusInternalServerError)
		return
	}

	// Restituisci "token"/ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(token.Token)

}
