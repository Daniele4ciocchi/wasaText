package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
)

// controlla se il token è valido e se si trova all'interno del db
// se non è valido restituisce un errore
// se è valido restituisce l'id dell'utente
func checkAuth(r *_router, req *http.Request) (int, error) {
	var token utils.Token
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return -1, fmt.Errorf("missing Authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return -1, fmt.Errorf("invalid Authorization header format")
	}
	token.Token = strings.TrimPrefix(authHeader, "Bearer ")
	if token.Token == "" {
		return -1, fmt.Errorf("empty token")
	}

	id, err := r.db.CheckToken(token.Token)
	if err != nil {
		return -1, fmt.Errorf("invalid token: %w", err)
	}
	return id, nil
}

// la funzione ritorna tutti gli utenti con cui è possibile inziare una conversazione
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	var users []utils.User

	users, err = rt.db.GetUsers()
	if err != nil {
		http.Error(w, "Errore nel recupero degli utenti", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// la funzione serve a ritornare un particolare utente
func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// auth control
	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	var user utils.User

	name := ps.ByName("userID")
	user, err = rt.db.GetUser(name)
	if err != nil {
		http.Error(w, "Errore nel recupero dell'utente", http.StatusInternalServerError) // da controllare se l'utente esiste, nel caso cambiare il tipo di errore
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
}

// con questa funzione possiamo impostare il nostro nome utente
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// auth control
	id, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	type Username struct {
		Username string `json:"username"`
	}
	var username Username
	err = json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, "Errore nella decodifica JSON", http.StatusBadRequest)
		return
	}

	err = rt.db.SetUsername(id, username.Username)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			http.Error(w, "Username già in uso", http.StatusConflict)
		} else {
			http.Error(w, "Errore nella modifica dell'username", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

// questa funzione ritorna le informazioni sull'utente loggato
func (rt *_router) getMe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// auth control
	id, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	user, err := rt.db.GetUserById(id)
	if err != nil {
		http.Error(w, "Errore nel recupero dell'utente", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
}

// con questa funzione possiamo impostare la nostra foto profilo
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	// Salva il file
	path, err := utils.SaveFile(user.Name, file)
	if err != nil {
		http.Error(w, "Errore nel salvataggio del file", http.StatusInternalServerError)
		return
	}

	// Aggiorna il path nel DB
	err = rt.db.SetUserPhoto(id, path)
	if err != nil {
		http.Error(w, "Errore nella modifica della foto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// la seguente funzione ritorna la foto di un determinato utente
func (rt *_router) getUserPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "image/jpeg")
	// auth control
	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "ID non valido", http.StatusBadRequest)
		return
	}

	// recupera il percorso della foto
	filePath, err := rt.db.GetUserPhoto(id)
	if err != nil {
		http.Error(w, "Errore nel recupero del percorso della foto", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(filePath)
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
