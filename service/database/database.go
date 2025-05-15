/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniele4ciocchi/wasaText/service/utils"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	//GetName() (string, error)
	//SetName(name string) error

	//user
	AddUser(name string, username string) error
	GetUser(name string) (utils.User, error)
	GetUsers() ([]utils.User, error)

	//conversation
	AddConversation(name string, isGroup bool) (int, error)
	AddUserConversation(userID int, convID int) error
	GetConversation(id int) (utils.Conversation, error)
	GetConversations(id int) ([]utils.Conversation, error)
	GetConversationByName(sender string, reciver string) (utils.Conversation, error)
	CheckExistingConversation(id1 int, id2 int) (int, error)

	//group
	AddGroup(name string) (utils.Group, error)
	GetGroupById(id int) (utils.Group, error)
	GetGroupByName(name string) (utils.Group, error)
	GetGroups(id int) ([]utils.Group, error)
	CheckExistingGroup(name string) (utils.Group, error)
	AddUserToGroup(userID int, groupID int) error
	RemoveUserFromGroup(userID int, groupID int) error
	GetGroupMembers(id int) ([]utils.User, error)
	LeaveGroup(userID int, groupID int) error

	//message
	AddMessage(senderID int, convID int, content string, repliedMessageID int) (int, error)
	GetMessage(id int) (utils.Message, error)
	GetMessages(convID int) ([]utils.Message, error)
	GetLastMessage(convID int) (utils.Message, error)

	//utils
	SetToken(id int, name string) error
	GetToken(id int) (string, error)
	GetUserFromToken(token string) (utils.User, error)
	CheckToken(token string) (int, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		sqlStmt := `
					CREATE TABLE users (
					    id INTEGER PRIMARY KEY AUTOINCREMENT,
					    name TEXT NOT NULL UNIQUE,
					    username TEXT NOT NULL UNIQUE
					);

					CREATE TABLE messages (
					    id INTEGER PRIMARY KEY AUTOINCREMENT,
					    conversation_id INTEGER NOT NULL,
					    sender_id INTEGER NOT NULL,
						replied_message_id INTEGER,
					    content TEXT NOT NULL,
					    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
					    FOREIGN KEY (conversation_id) REFERENCES conversations(id),
					    FOREIGN KEY (sender_id) REFERENCES users(id)
						FOREIGN KEY (replied_message_id) REFERENCES messages(id)
					);

					CREATE TABLE conversations (
					    id INTEGER PRIMARY KEY AUTOINCREMENT,
					    is_group BOOLEAN NOT NULL,
					    name TEXT NOT NULL
					);

					CREATE TABLE user_conversations (
					    id INTEGER PRIMARY KEY AUTOINCREMENT,
					    user_id INTEGER NOT NULL,
					    conversation_id INTEGER NOT NULL,
					    FOREIGN KEY (user_id) REFERENCES users(id),
					    FOREIGN KEY (conversation_id) REFERENCES conversations(id)
					);
					
					CREATE TABLE comments (
					    id INTEGER PRIMARY KEY AUTOINCREMENT,
					    message_id INTEGER NOT NULL,
					    user_id INTEGER NOT NULL,
					    content TEXT NOT NULL,
					    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
					    FOREIGN KEY (message_id) REFERENCES messages(id),
					    FOREIGN KEY (user_id) REFERENCES users(id)
					);

					CREATE TABLE tokens (
					    id INTEGER PRIMARY KEY AUTOINCREMENT,
						token TEXT NOT NULL,
						user_id INTEGER NOT NULL,
						FOREIGN KEY (user_id) REFERENCES users(id)
					);

					INSERT INTO users (name, username) VALUES ('admin', 'admin');

					`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
