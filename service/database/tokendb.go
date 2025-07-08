package database

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
)

// funzione per impostare il token di un nuovo utente
func (db *appdbimpl) SetToken(id int, name string) error {

	token := name + strconv.Itoa(id)
	_, err := db.c.Exec("INSERT OR REPLACE INTO tokens (token, user_id) VALUES (?, ?)", token, id)

	return err
}

// funzione per ottenere  il token di un utente
func (db *appdbimpl) GetToken(id int) (string, error) {
	var token string
	err := db.c.QueryRow("SELECT tokens.token FROM tokens WHERE user_id = ?", id).Scan(&token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return token, nil
}

// funzione per ottenere l'utente da un token
func (db *appdbimpl) GetUserFromToken(token string) (utils.User, error) {
	var id int
	var user utils.User
	err := db.c.QueryRow("SELECT user_id FROM tokens WHERE token = ?", token).Scan(&id)
	if err != nil {
		return utils.User{}, err
	}
	user.ID = id
	err = db.c.QueryRow("SELECT name, username FROM users WHERE id = ?", id).Scan(&user.Name, &user.Username)
	if err != nil {
		return utils.User{}, err
	}

	return user, err

}

// funzione per controllare se un determinato token esiste nel db
func (db *appdbimpl) CheckToken(token string) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT user_id FROM tokens WHERE token = ?", token).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
