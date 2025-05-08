package database

// funzione per aggiungere un utente al database
func (db *appdbimpl) AddUser(name string, username string) error {
	_, err := db.c.Exec("INSERT INTO users (name, username) VALUES (?, ?)", name, username)
	return err
}

func (db *appdbimpl) GetUser(name string) (int, string, string, error) {
	var id int
	var username string

	err := db.c.QueryRow("SELECT id, username FROM users WHERE name = ?", name).Scan(&id, &username)
	if err != nil {
		return 0, "", "", err
	}
	return id, name, username, nil
}
