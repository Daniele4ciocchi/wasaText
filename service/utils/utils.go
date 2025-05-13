package utils

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
