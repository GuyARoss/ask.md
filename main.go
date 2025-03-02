package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var lastResponseTime time.Time

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ask <command>")
		return
	}
	cmd := os.Args[1]
	if cmd == "watch" {
		watchFile("ask.md")
	} else {
		fmt.Println("Unknown command:", cmd)
	}
}

func watchFile(filename string) {
	var lastModTime time.Time
	for {
		fi, err := os.Stat(filename)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		modTime := fi.ModTime()
		if modTime.After(lastModTime) {
			lastModTime = modTime
			processFile(filename)
		}
		time.Sleep(1 * time.Second)
	}
}

func processFile(filename string) {
	// If a response was just appended, skip processing.
	if time.Since(lastResponseTime) < 2*time.Second {
		return
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}
	prompt := string(data)
	if len(strings.TrimSpace(prompt)) == 0 {
		// skipped
		return
	}

	answer, err := getAnswer(prompt)
	if err != nil {
		log.Println("Error getting answer:", err)
		return
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file for appending:", err)
		return
	}
	defer f.Close()

	// Append the response with a separator.
	_, err = f.WriteString("\n\n---\n\n" + answer)
	if err != nil {
		log.Println("Error writing answer to file:", err)
		return
	}
	// Record the time the LLM response was appended.
	lastResponseTime = time.Now()
}

func getAnswer(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}
	client := openai.NewClient(apiKey)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}
	return resp.Choices[0].Message.Content, nil
}
