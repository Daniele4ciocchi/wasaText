package api

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Conversation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID      int    `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}
