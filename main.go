package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func checkSentence(sentence string) (string, error) {
	prompt := `You are English Grammar teacher. For given sentence you will output whether it's gramatically correct or not. If not you will output only how it can be improved.
Keep in mind the original choice of words are try to preserve them as much as possible.  If the sentence is correct, only reply with "Correct"
Check this sentence. Don't repeat the original sentence.

"%s"`
	prompt = fmt.Sprintf(prompt, sentence)

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil

}

func main() {
	// Iterate over os.Args slice and print each argument
	arg := os.Args[1]

	b, err := os.ReadFile(arg)
	if err != nil {
		fmt.Println(err)
		return
	}
	contents := string(b)
	sentences := strings.Split(contents, "\n")[47:53]

	for _, sentence := range sentences {
		fmt.Printf("Checking sentence: %s\n", sentence)
		if sentence == "" {
			continue
		}
		result, err := checkSentence(sentence)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Result: %s\n", result)
	}

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

}
