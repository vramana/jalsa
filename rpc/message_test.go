package rpc

import (
	"testing"
)

func TestEncodeMessage(t *testing.T) {
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, world!",
	}

	message := EncodeMessage(data)

	expected := "Content-Length: 27\r\n\r\n{\"message\":\"Hello, world!\"}"

	if message != expected {
		t.Errorf("Expected %s, got %s", expected, message)
	}
}

func TestDecodeMessage(t *testing.T) {
	message := "Content-Length: 27\r\n\r\n{\"method\":\"hi\",\"message\":\"Hello, world!\"}"

	method, length, err := DecodeMessage([]byte(message))
	if len(length) != 27 {
		t.Errorf("Error decoding message: %v", err)
	}

	if method != "hi" {
		t.Errorf("Error decoding message: %v", err)
	}
}
