package database

import (
	"database/sql"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
)

// funzione per aggiungere un utente al database
func (db *appdbimpl) AddUser(name string, username string) error {
	_, err := db.c.Exec("INSERT INTO users (name, username) VALUES (?, ?)", name, username)
	if err != nil {
		return err
	}
	return nil
}

// funzione utilitaria che ritorna tutti gli utenti
func (db *appdbimpl) GetUsers() ([]utils.User, error) {
	rows, err := db.c.Query("SELECT id, name, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []utils.User
	for rows.Next() {
		var u utils.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// funzione che prende in input il nome di un utente e restituisce
// l'id, il nome e lo username
func (db *appdbimpl) GetUser(name string) (utils.User, error) {
	var user utils.User

	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT id, username FROM users WHERE name = ?)", name).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.User{}, nil
		}
		return utils.User{}, err
	}
	if !exists {
		return user, sql.ErrNoRows
	}
	user.Name = name
	err = db.c.QueryRow("SELECT id, username FROM users WHERE name = ?", name).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.User{}, nil
		}
		return utils.User{}, err
	}
	//controllo se l'utente è presente nel db
	if user.ID == 0 {
		return utils.User{}, nil
	}

	return user, nil
}
