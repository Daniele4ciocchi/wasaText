package utils

type Token struct {
	Token string `json:"Authorization"`
}

type Name struct {
	Name string `json:"name"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Conversation struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IsGroup bool   `json:"is_group"`
}

type Message struct {
	ID             int    `json:"id"`
	SenderID       int    `json:"sender_id"`
	ConversationID int    `json:"conversation_id"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
}
