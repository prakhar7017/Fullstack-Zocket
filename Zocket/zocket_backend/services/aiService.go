package services

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

var aiClient *openai.Client

func InitAI() {
	apiKey := os.Getenv("OPENAI_KEY")
	aiClient = openai.NewClient(apiKey)
}

func SuggestTask(prompt string) (string, error) {
	if aiClient == nil {
		log.Println("aiClient is nil! Make sure it is initialized properly.")
		return "", errors.New("AI client not initialized")
	}
	resp, err := aiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are an AI that suggests project tasks based on descriptions."},
				{Role: "user", Content: "Suggest tasks for: " + prompt},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func TaskBreakDown(prompt string) (string, error) {
	resp, err := aiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "You break down complex tasks into clear, actionable steps."},
				{Role: "user", Content: "Break down this task: " + prompt},
			},
		},
	)

	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func PrioritizeTasks(prompt string) (string, error) {

	resp, err := aiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are an AI that prioritizes tasks based on urgency and importance."},
				{Role: "user", Content: prompt},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
