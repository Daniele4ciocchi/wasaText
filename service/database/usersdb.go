package database

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) AddUser(name string, username string) error {
	_, err := db.c.Exec("INSERT INTO users (name, username) VALUES (?, ?)", name, username)
	return err
}

func (db *appdbimpl) GetAllUsernames() ([]string, error) {
	rows, err := db.c.Query("SELECT username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usernames, nil
}

func (db *appdbimpl) GetNameByUsername(username string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM users WHERE username = ?", username).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
