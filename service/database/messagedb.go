package database

import (
	"database/sql"
	"fmt"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
)

func (db *appdbimpl) AddMessage(senderID int, convID int, content string, repliedMessageID int) (int, error) {
	var messageID int

	var replied interface{}
	if repliedMessageID != 0 {
		replied = repliedMessageID
	} else {
		replied = sql.NullInt64{}
	}

	err := db.c.QueryRow("INSERT INTO messages (sender_id, conversation_id, content, replied_message_id ) VALUES (?, ?, ?, ?) RETURNING id", senderID, convID, content, replied).Scan(&messageID)
	if err != nil {
		fmt.Println(err)
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
	rows, err := db.c.Query("SELECT messages.id, messages.sender_id, users.name, messages.conversation_id, messages.replied_message_id, messages.content, messages.timestamp FROM messages JOIN users ON messages.sender_id = users.id WHERE conversation_id = ?", convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []utils.Message
	for rows.Next() {
		var message utils.Message
		var replied sql.NullInt64 // Utilizza NullInt64 per gestire il valore NULL

		if err := rows.Scan(&message.ID, &message.SenderID, &message.Sender, &message.ConversationID, &replied, &message.Content, &message.Timestamp); err != nil {
			return nil, err
		}

		// Verifica se replied è null o ha un valore
		if replied.Valid {
			// Usa una conversione esplicita da int64 a int
			message.RepliedMessageID = int(replied.Int64)
		} else {
			message.RepliedMessageID = 0
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (db *appdbimpl) GetLastMessage(convID int) (utils.Message, error) {
	var message utils.Message
	err := db.c.QueryRow("SELECT id, sender_id, conversation_id, content, timestamp FROM messages WHERE conversation_id = ? ORDER BY timestamp DESC LIMIT 1", convID).Scan(&message.ID, &message.SenderID, &message.ConversationID, &message.Content, &message.Timestamp)
	if err != nil {
		return utils.Message{}, err
	}

	return message, nil
}
