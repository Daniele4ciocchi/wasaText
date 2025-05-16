package utils

import (
	"io"
	"mime/multipart"
	"os"
)

type Token struct {
	Token string `json:"Authorization"`
}

type Name struct {
	Name string `json:"name"`
}

type User struct {
	ID       int    `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Group struct {
	ID      int    `json:"group_id"`
	Name    string `json:"name"`
	Members []User `json:"members"`
}

type Conversation struct {
	ID      int    `json:"conversation_id"`
	Name    string `json:"name"`
	IsGroup bool   `json:"is_group"`
}

type Message struct {
	ID               int    `json:"message_id"`
	SenderID         int    `json:"sender_id"`
	Sender           string `json:"sender"`
	RepliedMessageID int    `json:"replied_message_id"`
	ConversationID   int    `json:"conversation_id"`
	Content          string `json:"content"`
	Timestamp        string `json:"timestamp"`
}

func SaveFile(name string, fileContent multipart.File) (string, error) {
	// Define the file path
	basePath := "service/photos/"
	// Read the file content
	fullpath := basePath + name + ".jpg"
	content, err := io.ReadAll(fileContent)
	if err != nil {
		return "", err
	}
	// Create the file
	file, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// Write the content to the file
	_, err = file.Write(content)
	if err != nil {
		return "", err
	}
	return fullpath, nil
}
