package database

// funzione per aggiungere un utente al database
func (db *appdbimpl) AddUser(name string, username string) error {
	_, err := db.c.Exec("INSERT INTO users (name, username) VALUES (?, ?)", name, username)
	return err
}

// funzione utilitaria
func (db *appdbimpl) GetUsers() ([]string, error) {
	rows, err := db.c.Query("SELECT id, name, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []string
	for rows.Next() {
		var id int
		var name string
		var username string
		if err := rows.Scan(&id, &name, &username); err != nil {
			return nil, err
		}
		users = append(users, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// funzione che prende in input il nome di un utente e restituisce
// l'id, il nome e lo username
func (db *appdbimpl) GetUser(name string) (int, string, string, error) {
	var id int
	var username string

	err := db.c.QueryRow("SELECT id, username FROM users WHERE name = ?", name).Scan(&id, &username)
	if err != nil {
		return 0, "", "", err
	}
	return id, name, username, nil
}

func (db *appdbimpl) GetConversations(id int) ([]int, error) {
	rows, err := db.c.Query("SELECT conversation_id FROM user_conversations WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var convs []int
	for rows.Next() {
		var conv int
		if err := rows.Scan(&conv); err != nil {
			return nil, err
		}
		convs = append(convs, conv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return convs, nil
}
