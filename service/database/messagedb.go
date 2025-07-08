package database

import (
	"database/sql"
	"log"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
)

// funzione per aggiungere un messaggio al db
func (db *appdbimpl) AddMessage(senderID int, convID int, content string, repliedMessageID int, forwarded bool) (int, error) {
	var messageID int

	var replied interface{}
	if repliedMessageID != 0 {
		replied = repliedMessageID
	} else {
		replied = sql.NullInt64{}
	}

	err := db.c.QueryRow("INSERT INTO messages (sender_id, conversation_id, content, replied_message_id, forwarded) VALUES (?, ?, ?, ?, ?) RETURNING id", senderID, convID, content, replied, forwarded).Scan(&messageID)
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

// funzione per rimuovere un messaggio la db
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

// funzione per aggiungere il path di una foto al db come messaggio
func (db *appdbimpl) AddPhoto(senderID int, convID int, content string, repliedMessageID int, forwarded bool) (int, error) {
	var messageID int

	var replied interface{}
	if repliedMessageID != 0 {
		replied = repliedMessageID
	} else {
		replied = sql.NullInt64{}
	}

	err := db.c.QueryRow("INSERT INTO messages (sender_id, conversation_id, content, replied_message_id, photo, forwarded ) VALUES (?, ?, ?, ?, true, ?) RETURNING id", senderID, convID, content, replied, forwarded).Scan(&messageID)
	if err != nil {
		return 0, err
	}

	err = db.SetSendedMessage(messageID)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}

// funzione per ottenere un determinato messaggio contenuto nel db
func (db *appdbimpl) GetMessage(id int) (utils.Message, error) {
	var message utils.Message
	err := db.c.QueryRow("SELECT id, sender_id, conversation_id, content, photo, forwarded, timestamp FROM messages WHERE id = ?", id).Scan(&message.ID, &message.SenderID, &message.ConversationID, &message.Content, &message.Photo, &message.Forwarded, &message.Timestamp)
	if err != nil {
		return utils.Message{}, err
	}
	return message, nil
}

// funzione per ottenere tutti i messaggi di una conversazione
func (db *appdbimpl) GetMessages(convID int) ([]utils.Message, error) {
	rows, err := db.c.Query("SELECT messages.id, messages.sender_id, users.name, messages.conversation_id, messages.replied_message_id, messages.content, messages.timestamp, messages.photo, messages.forwarded FROM messages JOIN users ON messages.sender_id = users.id WHERE conversation_id = ?", convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []utils.Message
	var messageIDsToUpdate []int // Slice to store IDs of messages that need status update

	for rows.Next() {
		var message utils.Message
		var replied sql.NullInt64 // Utilizza NullInt64 per gestire il valore NULL

		if err := rows.Scan(&message.ID, &message.SenderID, &message.Sender, &message.ConversationID, &replied, &message.Content, &message.Timestamp, &message.Photo, &message.Forwarded); err != nil {
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
		messageIDsToUpdate = append(messageIDsToUpdate, message.ID) // Collect message ID for later update
	}

	// Check for any errors that occurred during rows.Next()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Now, after all messages have been read and rows are closed,
	// update the changed_status for each collected message ID.
	for _, messageID := range messageIDsToUpdate {
		_, err := db.c.Exec("UPDATE messages SET changed_status = FALSE WHERE id = ?", messageID)
		if err != nil {
			// Log the error. You might not want to return an error here
			// as the messages themselves have already been successfully retrieved.
			log.Printf("Error setting changed_status to FALSE for message %d: %v\n", messageID, err)
		}
	}

	return messages, nil
}

// funzione per ottenere lo stato di un messaggio (letto, consegnato, inviato)
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

// funzione per assegnare lo stato di invio ad un messaggio
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

// funzione per impostare che un messaggio è arrivato
func (db *appdbimpl) SetArrivedMessage(userID int, messageID int) error {
	_, err := db.c.Exec("UPDATE views SET status = ? WHERE user_id = ? AND message_id = ?", 1, userID, messageID)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("UPDATE messages SET changed_status = ? WHERE id = ?", true, messageID)
	if err != nil {
		return err
	}

	return nil
}

// funzione per impostare la lettura di un messaggio
func (db *appdbimpl) SetViewedMessage(userID int, messageID int) error {
	_, err := db.c.Exec("UPDATE views SET status = ? WHERE user_id = ? AND message_id = ?", 2, userID, messageID)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("UPDATE messages SET changed_status = ? WHERE id = ?", true, messageID)
	if err != nil {
		return err
	}
	return nil
}

// funzione per ottenere l'ultimo messaggio di una conversaione
func (db *appdbimpl) GetLastMessage(convID int) (utils.Message, error) {
	var message utils.Message
	err := db.c.QueryRow("SELECT id, sender_id, conversation_id, content, photo, forwarded, timestamp FROM messages WHERE conversation_id = ? ORDER BY timestamp DESC LIMIT 1", convID).Scan(&message.ID, &message.SenderID, &message.ConversationID, &message.Content, &message.Photo, &message.Forwarded, &message.Timestamp)
	if err != nil {
		return utils.Message{}, err
	}

	return message, nil
}

// funzione per ottenere nuovi messaggi di una conversazione (non letti)
func (db *appdbimpl) GetNewMessages(userID int) ([]utils.Message, error) {
	rows, err := db.c.Query("SELECT m.id, m.sender_id, m.conversation_id, m.content, m.forwarded, m.timestamp FROM messages m JOIN views v ON m.id = v.message_id WHERE v.user_id = ? AND (v.status = 0 OR m.changed_status = 1) ", userID)
	if err != nil {
		// print error
		log.Println("Error getting new messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []utils.Message
	for rows.Next() {
		var message utils.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ConversationID, &message.Content, &message.Forwarded, &message.Timestamp); err != nil {
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
