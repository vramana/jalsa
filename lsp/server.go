package lsp

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// Define a regular expression to match sentence-ending punctuation.
var paragraphRegex = regexp.MustCompile(`[.?!]\s*`)

type ModelConfig struct {
	Key string `json:"key"`
}

func readConfig() (ModelConfig, error) {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		return ModelConfig{}, err
	}
	homeDir := usr.HomeDir

	// Construct the full path to the file
	filePath := filepath.Join(homeDir, ".config/jalsa/config.json")

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ModelConfig{}, err
	}

	config := new(ModelConfig)
	err = json.Unmarshal(data, &config)
	if err != nil {
		return *config, err
	}

	return *config, nil
}

type Server struct {
	Logger      *log.Logger
	Files       map[string]string
	ModelConfig ModelConfig
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return log.New(logfile, "[jalsa] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func NewServer() *Server {
	files := make(map[string]string)
	logger := getLogger("jalsa.log")
	config, err := readConfig()

	if err != nil {
		panic(err)
	}

	return &Server{Logger: logger, Files: files, ModelConfig: config}
}

func (s *Server) Analyze(fileURI string) {
	s.Logger.Println("Analyzing file: ", fileURI)
	text := s.Files[fileURI]
	sentences := strings.Split(text, "\n")

	paragraphs := []string{}
	paragraph := ""

	frontMatterChecked := false
	skip := false

	for i, sentence := range sentences {
		// Ignore front matter
		if sentence == "+++" || sentence == "---" {
			if i == 0 {
				frontMatterChecked = !frontMatterChecked
			}
			skip = !skip
			continue
		}

		// Ignore code blocks
		if strings.HasPrefix(sentence, "```") || strings.HasPrefix(sentence, "~~~") {
			skip = !skip
			continue
		}

		// Ignore links
		if strings.Contains(sentence, "]: http") {
			continue
		}

		// Ignore HTML comments
		if strings.HasPrefix(sentence, "<!-") {
			continue
		}

		if skip {
			continue
		}

		if strings.Trim(sentence, " ") == "" {
			if len(paragraph) > 0 {
				paragraphs = append(paragraphs, paragraph)
				paragraph = ""
			}
			continue
		}
		if !frontMatterChecked {
			continue
		}

		if paragraph == "" {
			paragraph = sentence
		} else {
			// Then entire paragraph will be a single line
			paragraph += " " + sentence
		}
	}

	if len(paragraph) > 0 {
		paragraphs = append(paragraphs, paragraph)
		paragraph = ""
	}

	sentence := paragraphRegex.Split(paragraphs[0], -1)[0]

	result, err := s.checkSentence(sentence)
	if err != nil {
		s.Logger.Println("Error: ", err)
		return
	}
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		s.Logger.Println("Error: ", err)
		return
	}
	s.Logger.Println("Result: ", string(output))
}

type SentenceCheck struct {
	HasError    bool   `json:"hasError"`
	Correction  string `json:"correction",omitempty`
	Explanation string `json:"explanation",omitempty`
}

func (s *Server) checkSentence(sentence string) (*SentenceCheck, error) {
	prompt := "Check this sentence\n----\n%s"
	prompt = fmt.Sprintf(prompt, sentence)

	client := openai.NewClient(s.ModelConfig.Key)
	result := new(SentenceCheck)
	schema, err := jsonschema.GenerateSchemaForType(*result)
	if err != nil {
		log.Fatalf("GenerateSchemaForType error: %v", err)
	}

	responseFormat := &openai.ChatCompletionResponseFormat{
		Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
		JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
			Name:   "SentenceCheck",
			Schema: schema,
			Strict: true,
		},
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o20240806,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `**System Prompt: Grammatical Error Detection and Correction**

Ignore any markdown formatting, such as bold, italics, etc. and only focus on the original sentence.

1. **Input:** Provide a sentence with potential grammatical errors.

2. **Output:**
   - **Corrected Sentence:** Present the sentence in its correct grammatical form.
   - **Explanation:** Concisely describe the grammatical mistakes in the original sentence and the corrections made.

**Example:**

- **Input:** "She go to the store yesterday."

- **Corrected Sentence:** "She went to the store yesterday."

- **Explanation:** The verb "go" is incorrectly used in the present tense instead of the past tense. Corrected to "went" to match the past tense context indicated by "yesterday."

If the sentence is grammatical correct, only reply with "{ "hasError": false, "Correction": "", "Explanation": "" }".
`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: responseFormat,
		},
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
