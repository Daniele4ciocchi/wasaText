package database

import "strconv"

func (db *appdbimpl) SetToken(id int, name string) error {

	token := name + strconv.Itoa(id)
	_, err := db.c.Exec("INSERT INTO tokens (token, user_id) VALUES (?, ?)", token, id)
	return err
}

func (db *appdbimpl) GetToken(name string) (string, error) {
	var token string
	err := db.c.QueryRow("SELECT tokens.token FROM tokens JOIN users ON tokens.users_id = users.id WHERE users.name = ?", name).Scan(&token)
	return token, err
}

func (db *appdbimpl) CheckToken(token string) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT user_id FROM tokens WHERE token = ?", token).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
