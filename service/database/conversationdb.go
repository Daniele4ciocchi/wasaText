package database

import (
	"github.com/Daniele4ciocchi/wasaText/service/utils"
)

// AddConversation adds a new conversation to the database.
func (db *appdbimpl) AddConversation(name string, isGroup bool) (int, error) {
	var convID int
	err := db.c.QueryRow("INSERT INTO conversations (name, is_group) VALUES (?, ?) RETURNING id", name, isGroup).Scan(&convID)
	if err != nil {
		return 0, err
	}
	return convID, nil
}

// AddUserConversation adds a user to a conversation.
func (db *appdbimpl) AddUserConversation(userID int, convID int) error {
	_, err := db.c.Exec("INSERT INTO user_conversations (user_id, conversation_id) VALUES (?, ?)", userID, convID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetConversation(id int) (utils.Conversation, error) {
	var conv utils.Conversation
	err := db.c.QueryRow("SELECT id, name, is_group FROM conversations WHERE id = ?", id).Scan(&conv.ID, &conv.Name, &conv.IsGroup)
	if err != nil {
		return utils.Conversation{}, err
	}

	return conv, nil

}

func (db *appdbimpl) GetConversations(id int) ([]utils.Conversation, error) {
	rows, err := db.c.Query("SELECT conversation_id FROM user_conversations WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []utils.Conversation
	for rows.Next() {
		var convID int
		if err := rows.Scan(&convID); err != nil {
			return nil, err
		}
		conv, err := db.GetConversation(convID)
		if err != nil {
			return nil, err
		}
		convs = append(convs, conv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return convs, nil
}
