package database

import (
	"github.com/Daniele4ciocchi/wasaText/service/utils"
)

func (db *appdbimpl) GetGroupById(id int) (utils.Group, error) {
	var group utils.Group
	err := db.c.QueryRow("SELECT id, name FROM conversations WHERE id = ? AND is_group = 1", id).Scan(&group.ID, &group.Name)
	if err != nil {
		return utils.Group{}, err
	}
	return group, nil
}

func (db *appdbimpl) GetGroupByName(name string) (utils.Group, error) {
	var group utils.Group
	err := db.c.QueryRow("SELECT id, name FROM conversations WHERE name = ? AND is_group = 1", name).Scan(&group.ID, &group.Name)
	if err != nil {
		return utils.Group{}, err
	}
	return group, nil
}

// GetGroups returns all groups from user id
func (db *appdbimpl) GetGroups(id int) ([]utils.Group, error) {
	rows, err := db.c.Query("SELECT c.id, c.name FROM user_conversations uc JOIN conversations c ON uc.conversation_id = c.id WHERE uc.user_id = ? AND c.is_group = 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []utils.Group
	for rows.Next() {
		var group utils.Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// CheckExistingGroup controlla se esiste un gruppo, se esiste torna l'id
func (db *appdbimpl) CheckExistingGroup(name string) (utils.Group, error) {
	var group utils.Group
	err := db.c.QueryRow(`SELECT id FROM conversations WHERE conversations.name = ? AND is_group = 1`, name).Scan(&group.ID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			group.ID = 0
			return group, nil
		}
		return utils.Group{}, err
	}
	group.Name = name
	return group, nil
}

func (db *appdbimpl) AddGroup(name string) (utils.Group, error) {
	// Add group to conversations table
	var group utils.Group
	res, err := db.c.Exec("INSERT INTO conversations (name, is_group, photoPath) VALUES (?, 1, 'service/photos/default.jpg')", name)
	if err != nil {
		group.ID = 0
		return group, err
	}

	// Get the last inserted ID
	lastID, err := res.LastInsertId()
	if err != nil {
		group.ID = 0
		return group, err
	}
	group.ID = int(lastID)

	return group, nil
}

func (db *appdbimpl) CheckUserInGroup(userID int, groupID int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM user_conversations WHERE user_id = ? AND conversation_id = ?)", userID, groupID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *appdbimpl) AddUserToGroup(userID int, groupID int) error {
	_, err := db.c.Exec("INSERT INTO user_conversations (user_id, conversation_id) VALUES (?, ?)", userID, groupID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveUserFromGroup(userID int, groupID int) error {
	_, err := db.c.Exec("DELETE FROM user_conversations WHERE user_id = ? AND conversation_id = ?", userID, groupID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetGroupMembers(groupID int) ([]utils.User, error) {
	rows, err := db.c.Query("SELECT u.id, u.name, u.username FROM user_conversations uc JOIN users u ON uc.user_id = u.id WHERE uc.conversation_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []utils.User
	for rows.Next() {
		var member utils.User
		if err := rows.Scan(&member.ID, &member.Name, &member.Username); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (db *appdbimpl) LeaveGroup(userID int, groupID int) error {
	_, err := db.c.Exec("DELETE FROM user_conversations WHERE user_id = ? AND conversation_id = ?", userID, groupID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetGroupName(id int, name string) error {
	_, err := db.c.Exec("UPDATE conversations SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetGroupPhoto(id int, path string) error {
	_, err := db.c.Exec("UPDATE conversations SET photoPath = ? WHERE id = ?", path, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetGroupPhoto(id int) (string, error) {
	var path string
	err := db.c.QueryRow("SELECT photoPath FROM conversations WHERE id = ? AND is_group + true", id).Scan(&path)
	if err != nil {
		return "", err
	}
	return path, nil
}
