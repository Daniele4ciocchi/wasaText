package database

import (
	"database/sql"
	"log"

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
		return 0, err
	}

	err = db.SetSendedMessage(messageID)
	if err != nil {
		return 0, err
	}

	err = db.SetArrivedMessage(senderID, messageID)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}

func (db *appdbimpl) RemoveMessage(messageID int) error {
	// Rimuove il messaggio dalla tabella views
	_, err := db.c.Exec("DELETE FROM views WHERE message_id = ?", messageID)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM messages WHERE id = ?", messageID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) AddPhoto(senderID int, convID int, content string, repliedMessageID int) (int, error) {
	var messageID int

	var replied interface{}
	if repliedMessageID != 0 {
		replied = repliedMessageID
	} else {
		replied = sql.NullInt64{}
	}

	err := db.c.QueryRow("INSERT INTO messages (sender_id, conversation_id, content, replied_message_id, photo ) VALUES (?, ?, ?, ?, true) RETURNING id", senderID, convID, content, replied).Scan(&messageID)
	if err != nil {
		return 0, err
	}

	err = db.SetSendedMessage(messageID)
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
	rows, err := db.c.Query("SELECT messages.id, messages.sender_id, users.name, messages.conversation_id, messages.replied_message_id, messages.content, messages.timestamp, messages.photo FROM messages JOIN users ON messages.sender_id = users.id WHERE conversation_id = ?", convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []utils.Message
	for rows.Next() {
		var message utils.Message
		var replied sql.NullInt64 // Utilizza NullInt64 per gestire il valore NULL

		if err := rows.Scan(&message.ID, &message.SenderID, &message.Sender, &message.ConversationID, &replied, &message.Content, &message.Timestamp, &message.Photo); err != nil {
			return nil, err
		}

		// Verifica se replied è null o ha un valore
		if replied.Valid {
			// Usa una conversione esplicita da int64 a int
			message.RepliedMessageID = int(replied.Int64)
		} else {
			message.RepliedMessageID = 0
		}
		message.Status, err = db.GetMessageStatus(message.ID)
		if err != nil {
			log.Println(err)
		}
		messages = append(messages, message)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (db *appdbimpl) GetMessageStatus(messageID int) (int, error) {
	var status int

	var mess utils.Message

	mess, err := db.GetMessage(messageID)
	if err != nil {
		return 0, err
	}

	var users []utils.User
	users, err = db.GetConversationMembers(mess.ConversationID)
	if err != nil {
		return 0, err
	}

	var sended, arrived, viewed int
	for _, user := range users {
		err := db.c.QueryRow("SELECT status FROM views WHERE message_id = ? AND user_id = ?", messageID, user.ID).Scan(&status)
		if err != nil {
			return 0, err
		}
		if status == 0 {
			sended++
		}
		if status == 1 {
			sended++
			arrived++
		}
		if status == 2 {
			sended++
			arrived++
			viewed++
		}
	}
	TotUsers := len(users)
	if TotUsers == viewed {
		return 2, nil
	}
	if TotUsers == arrived {
		return 1, nil
	}
	if TotUsers == sended {
		return 0, nil
	}
	return 0, nil
}

func (db *appdbimpl) SetSendedMessage(messageID int) error {

	var mess utils.Message
	mess, err := db.GetMessage(messageID)
	if err != nil {
		return err
	}

	var users []utils.User
	users, err = db.GetConversationMembers(mess.ConversationID)
	if err != nil {
		return err
	}

	for _, user := range users {
		_, err := db.c.Exec("INSERT INTO views (message_id, user_id, status) VALUES (?, ?, ?) ", messageID, user.ID, 0)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (db *appdbimpl) SetArrivedMessage(userID int, messageID int) error {
	_, err := db.c.Exec("UPDATE views SET status = ? WHERE user_id = ? AND message_id = ?", 1, userID, messageID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetViewedMessage(userID int, messageID int) error {
	_, err := db.c.Exec("UPDATE views SET status = ? WHERE user_id = ? AND message_id = ?", 2, userID, messageID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetLastMessage(convID int) (utils.Message, error) {
	var message utils.Message
	err := db.c.QueryRow("SELECT id, sender_id, conversation_id, content, timestamp FROM messages WHERE conversation_id = ? ORDER BY timestamp DESC LIMIT 1", convID).Scan(&message.ID, &message.SenderID, &message.ConversationID, &message.Content, &message.Timestamp)
	if err != nil {
		return utils.Message{}, err
	}

	return message, nil
}

// ritorna una lista di messaggi non arrivati e li segna come arrivati
func (db *appdbimpl) GetNewMessages(userID int) ([]utils.Message, error) {
	rows, err := db.c.Query("SELECT m.id, m.sender_id, m.conversation_id, m.content, m.timestamp FROM messages m JOIN views v ON m.id = v.message_id WHERE v.user_id = ? AND v.status = 0 ", userID)
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
		message.Status, err = db.GetMessageStatus(message.ID)
		if err != nil {
			return nil, err
		}
		var user utils.User
		user, err = db.GetUserById(message.SenderID)
		if err != nil {
			return nil, err
		}
		message.Sender = user.Name
		messages = append(messages, message)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
