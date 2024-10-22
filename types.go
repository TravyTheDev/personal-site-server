package main

type ChatMessage struct {
	UserID    string `json:"userID"`
	Name      string `json:"name"`
	Body      string `json:"body"`
	MessageID string `json:"messageID"`
}

type Message struct {
	Message string `json:"msg"`
}

type MessageWithSubtraction struct {
	Message    string `json:"message"`
	FirstNum   int    `json:"firstNum"`
	SecondNum  int    `json:"secondNum"`
	Difference int    `json:"difference"`
}
