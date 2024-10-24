package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
}

func EncodeMessage(data any) string {
	content, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(message []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(message, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("no content length found")
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMessage BaseMessage
	err = json.Unmarshal(content, &baseMessage)
	if err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, nil
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	messageLength := len(header) + contentLength + 4

	return messageLength, data[:messageLength], nil
}
