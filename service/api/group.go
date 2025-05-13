package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := checkAuth(rt, r)
	if err != nil {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		http.Error(w, "ID del gruppo non valido", http.StatusBadRequest)
		return
	}

	var group utils.Group
	group, err = rt.db.GetGroupById(groupID)
	println(groupID)
	if err != nil {
		http.Error(w, "Errore nel recupero del gruppo", http.StatusInternalServerError)
		return
	}
	var users []utils.User
	users, err = rt.db.GetGroupMembers(groupID)
	if err != nil {
		http.Error(w, "Errore nel recupero degli utenti", http.StatusInternalServerError)
		return
	}
	group.Members = users

	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	//auth control
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
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	type newGroup struct {
		Name    string `json:"name"`
		Members []int  `json:"members"`
	}

	// decodifico il JSON
	var group newGroup
	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, "Errore nella decodifica JSON", http.StatusBadRequest)
		return
	}

	for _, member := range group.Members {
		if member == user.ID {
			http.Error(w, "Non puoi aggiungere te stesso al gruppo", http.StatusBadRequest)
			return
		}
	}

	//controllo se il gruppo esiste già
	existingGroup, err := rt.db.CheckExistingGroup(group.Name)
	if err != nil {
		http.Error(w, "Errore durante il recupero del gruppo", http.StatusInternalServerError)
		return
	}
	if existingGroup.ID != 0 {
		http.Error(w, "Il gruppo esiste già", http.StatusConflict)
		return
	}

	//creo il gruppo
	var createdGroup utils.Group
	createdGroup, err = rt.db.AddGroup(group.Name)
	if err != nil {
		http.Error(w, "Errore durante la creazione del gruppo", http.StatusInternalServerError)
		return
	}

	//aggiungo gli utenti al gruppo
	err = rt.db.AddUserToGroup(user.ID, createdGroup.ID)
	if err != nil {
		http.Error(w, "Errore durante l'aggiunta dell'utente al gruppo", http.StatusInternalServerError)
		return
	}
	for _, member := range group.Members {
		println(member)
		err = rt.db.AddUserToGroup(member, createdGroup.ID)
		if err != nil {
			http.Error(w, "Errore durante l'aggiunta dell'utente al gruppo", http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(createdGroup); err != nil {
		http.Error(w, "Errore nella codifica JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	//auth control
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
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		http.Error(w, "ID del gruppo non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.LeaveGroup(user.ID, groupID)
	if err != nil {
		http.Error(w, "Errore durante l'uscita dal gruppo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
