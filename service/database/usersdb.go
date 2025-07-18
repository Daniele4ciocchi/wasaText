package database

import (
	"database/sql"
	"errors"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
)

// funzione per aggiungere un utente al database
func (db *appdbimpl) AddUser(name string, username string) error {
	_, err := db.c.Exec("INSERT INTO users (name, username, photoPath) VALUES (?, ?, 'service/photos/default.jpg')", name, username)
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

// funzione per trovare un utente dal suo id
func (db *appdbimpl) GetUserById(id int) (utils.User, error) {
	var user utils.User
	err := db.c.QueryRow("SELECT id, name, username FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, sql.ErrNoRows
		}
		return user, err
	}
	return user, nil
}

// funzione che prende in input il nome di un utente e restituisce
// l'id, il nome e lo username
func (db *appdbimpl) GetUser(name string) (utils.User, error) {
	var user utils.User

	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT id, username FROM users WHERE name = ?)", name).Scan(&exists)
	if err != nil {
		return user, err
	}
	if !exists {
		return user, sql.ErrNoRows
	}
	user.Name = name
	err = db.c.QueryRow("SELECT id, username FROM users WHERE name = ?", name).Scan(&user.ID, &user.Username)
	if err != nil {
		return utils.User{}, err
	}

	return user, nil
}

// funzione per impostare lo username di un utente
func (db *appdbimpl) SetUsername(id int, username string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", username, id)
	if err != nil {
		return err
	}
	return nil
}

// funzione per impostare la foto di un utente
func (db *appdbimpl) SetUserPhoto(id int, path string) error {

	_, err := db.c.Exec("UPDATE users SET photoPath = ? WHERE id = ?", path, id)
	if err != nil {
		return err
	}
	return nil
}

// funzione per ottenere il path della foto di un utente
func (db *appdbimpl) GetUserPhoto(id int) (string, error) {
	var path string
	err := db.c.QueryRow("SELECT photoPath FROM users WHERE id = ?", id).Scan(&path)
	if err != nil {
		return "", err
	}
	return path, nil
}
