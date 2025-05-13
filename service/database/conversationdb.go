package database

import (
	"fmt"

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

// controlla tramite due nomi se la conversazione esiste già e ritorna l'id della conversazione
func (db *appdbimpl) CheckExistingConversation(id1 int, id2 int) (int, error) {
	var convID int
	err := db.c.QueryRow(`SELECT uc1.conversation_id
							FROM user_conversations uc1 JOIN user_conversations uc2 ON uc1.conversation_id = uc2.conversation_id
							WHERE uc1.user_id = ? AND uc2.user_id = ?
							LIMIT 1`, id1, id2).Scan(&convID)
	if err != nil {
		fmt.Println(err)
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

func (db *appdbimpl) GetConversationByName(sender string, reciver string) (utils.Conversation, error) {
	var conv utils.Conversation
	err := db.c.QueryRow("SELECT id, name, is_group FROM conversations WHERE (name = ? OR name = ?) AND is_group = false", sender+reciver, reciver+sender).Scan(&conv.ID, &conv.Name, &conv.IsGroup)
	if err != nil {
		return utils.Conversation{}, err
	}
	return conv, nil
}
