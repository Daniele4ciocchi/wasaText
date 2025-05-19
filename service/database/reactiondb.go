package database

import "github.com/Daniele4ciocchi/wasaText/service/utils"

func (db *appdbimpl) AddReaction(userID int, messageID int, content string) error {
	_, err := db.c.Exec("INSERT INTO reactions (user_id, message_id, content) VALUES (?, ?, ?)", userID, messageID, content)
	return err
}

func (db *appdbimpl) RemoveReaction(reactionID int) error {
	_, err := db.c.Exec("DELETE FROM reactions WHERE id = ?", reactionID)
	return err
}

func (db *appdbimpl) GetReactions(messageID int) ([]utils.Reaction, error) {
	rows, err := db.c.Query("SELECT id, user_id, content FROM reactions WHERE message_id = ?", messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []utils.Reaction
	for rows.Next() {
		var reaction utils.Reaction
		if err := rows.Scan(&reaction.ID, &reaction.UserID, &reaction.Content); err != nil {
			return nil, err
		}
		reaction.MessageID = messageID
		reactions = append(reactions, reaction)
	}

	return reactions, nil
}
