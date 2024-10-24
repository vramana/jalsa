package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"log"
	// "database/sql"
	// "bytes"
	"fmt"
	"io"
	// "log"
	"os"
	// "strings"
	// "time"
	// _ "github.com/mattn/go-sqlite3"
	"github.com/sashabaranov/go-openai"

	"jalsa/lsp"
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

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// func test() {
// 	// Iterate over os.Args slice and print each argument
//
// 	db, err := sql.Open("sqlite3", "./test.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
//
// 	arg := os.Args[1]
// 	b, err := os.ReadFile(arg)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	fmt.Printf("Reading file: %s\n", arg)
//
// 	contents := string(b)
// 	sentences := strings.Split(contents, "\n")
//
// 	frontMatter := false
//
// 	for _, sentence := range sentences {
// 		if sentence == "+++" {
// 			frontMatter = !frontMatter
// 			continue
// 		}
// 		if sentence == "" {
// 			continue
// 		}
// 		fmt.Printf("Checking sentence: %s\n", sentence)
// 		h := hash(sentence)
// 		var result string
// 		err = db.QueryRow("SELECT correction FROM sentences WHERE sentence_hash = ?", h).Scan(&result)
//
// 		if err != nil && err != sql.ErrNoRows {
// 			fmt.Printf("Error: %v\n", err)
// 			continue
// 		}
//
// 		if result != "" {
// 			fmt.Printf("Result: %s\n", result)
// 			continue
// 		}
//
// 		fmt.Printf("Not found in cache\n")
//
// 		result, err = checkSentence(sentence)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			continue
// 		}
// 		fmt.Printf("Result: %s\n", result)
//
// 		time.Sleep(100 * time.Millisecond)
//
// 		_, err = db.Exec("INSERT INTO sentences (sentence_hash, sentence, correction) VALUES (?, ?, ?)", hash(sentence), sentence, result)
// 	}
// }

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(Split)

	logger := getLogger("jalsa.log")

	logger.Println("Starting...")

	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, msg, err := DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error %s", err)
			continue
		}
		handleMessage(msg, writer, method, logger)
	}
}

func handleMessage(msg []byte, writer io.Writer, method string, logger *log.Logger) {
	logger.Printf("Method: %s", method)
	logger.Println(string(msg))

	switch method {
	case "initialize":
		logger.Println(string(msg))
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return log.New(logfile, "[jalsa] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func writeMessage(writer io.Writer, message any) {
	data := EncodeMessage(message)
	writer.Write([]byte(data))
}
