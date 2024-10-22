package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var history []*genai.Content

type ChatBot struct {
	ctx context.Context
}

func NewChatBot() *ChatBot {
	return &ChatBot{}
}

func (c *ChatBot) Chat(w http.ResponseWriter, r *http.Request) {
	c.ctx = context.Background()
	apiKey := os.Getenv("API_KEY")

	var input MessageWithSubtraction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("HERE", err)
	}
	if input.FirstNum-input.SecondNum != input.Difference {
		return
	}
	client, err := genai.NewClient(c.ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Println("HEY", err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()

	cs.History = history

	res, err := cs.SendMessage(c.ctx, genai.Text(input.Message))
	if err != nil {
		log.Println("HELLO", err)
	}

	userInput := genai.Content{
		Parts: []genai.Part{
			genai.Text(input.Message),
		},
		Role: "user",
	}

	history = append(history, &userInput)

	botRes := printResponse(res)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(botRes); err != nil {
		http.Error(w, "error getting response from bot", http.StatusInternalServerError)
		return
	}
}

func printResponse(resp *genai.GenerateContentResponse) []string {
	finalString := []string{}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				str := fmt.Sprint(part)
				botRes := genai.Content{
					Parts: []genai.Part{
						genai.Text(str),
					},
					Role: "model",
				}
				history = append(history, &botRes)
				finalString = append(finalString, str)
			}
		}
	}
	return finalString
}
