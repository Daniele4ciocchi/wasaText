package database

import "github.com/Daniele4ciocchi/wasaText/service/utils"

func (db *appdbimpl) AddMessage(senderID int, convID int, content string) (int, error) {
	var messageID int
	err := db.c.QueryRow("INSERT INTO messages (sender_id, conversation_id, content) VALUES (?, ?, ?) RETURNING id", senderID, convID, content).Scan(&messageID)
	if err != nil {
		return 0, err
	}
	return messageID, nil
}

func (db *appdbimpl) GetMessage(id int) (utils.Message, error) {
	var message utils.Message
	err := db.c.QueryRow("SELECT id, sender_id, conversation_id, content, timestamp FROM messages WHERE id = ?", id).Scan(&message.ID, &message.SenderID, &message.ConversationID, &message.Content, &message.Timestamp)
	if err != nil {
		return utils.Message{}, err
	}
	return message, nil
}

func (db *appdbimpl) GetMessages(convID int) ([]utils.Message, error) {
	rows, err := db.c.Query("SELECT id, sender_id, conversation_id, content, timestamp FROM messages WHERE conversation_id = ?", convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []utils.Message
	for rows.Next() {
		var message utils.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ConversationID, &message.Content, &message.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
