package database

import "database/sql"

func (db *appdbimpl) AddConversation(name string, isGroup bool) (int, error) {
	var convID int
	err := db.c.QueryRow("INSERT INTO conversations (name, is_group) VALUES (?, ?) RETURNING id", name, isGroup).Scan(&convID)
	if err != nil {
		return 0, err
	}
	return convID, nil
}

func (db *appdbimpl) GetConversation(id int) (int, error) {
	var conv int
	err := db.c.QueryRow("SELECT id FROM conversations WHERE id = ?", id).Scan(&conv)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return conv, nil
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
