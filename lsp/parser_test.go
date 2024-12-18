package lsp

import (
	"testing"
)

type ParserTest struct {
	Text     string
	Expected []Sentence
}

func TestParse(t *testing.T) {
	tests := []ParserTest{
		ParserTest{
			Text: `Hello, world! This is a test.`,
			Expected: []Sentence{
				Sentence{"Hello, world", Range{Position{0, 0}, Position{0, 12}}},
				Sentence{"This is a test", Range{Position{0, 14}, Position{0, 28}}},
			},
		},
		ParserTest{
			Text: "Hello, world! This is a\ntest. This is another test.",
			Expected: []Sentence{
				Sentence{"Hello, world", Range{Position{0, 0}, Position{0, 12}}},
				Sentence{"This is a test", Range{Position{0, 14}, Position{1, 4}}},
				Sentence{"This is another test", Range{Position{1, 6}, Position{1, 26}}},
			},
		},
		ParserTest{
			Text: "- Hello world. This is\n  a sentence\n- This is a test",
			Expected: []Sentence{
				Sentence{"- Hello world", Range{Position{0, 0}, Position{0, 13}}},
				Sentence{"This is   a sentence", Range{Position{0, 15}, Position{1, 12}}},
				Sentence{"- This is a test", Range{Position{2, 0}, Position{2, 16}}},
			},
		},
		ParserTest{
			Text: "Hello world\n  \nThis is a sentence\n\n- This is a test",
			Expected: []Sentence{
				Sentence{"Hello world", Range{Position{0, 0}, Position{0, 11}}},
				Sentence{"This is a sentence", Range{Position{2, 0}, Position{2, 18}}},
				Sentence{"- This is a test", Range{Position{4, 0}, Position{4, 16}}},
			},
		},
		ParserTest{
			Text: `I am trying to write a simple Redis client in Rust. As part of that exercise, I have read the source code of
Node.js Redis Client. I want to jot down some notes on how it works.`,
			Expected: []Sentence{
				Sentence{"I am trying to write a simple Redis client in Rust", Range{Position{0, 0}, Position{0, 50}}},
				Sentence{"As part of that exercise, I have read the source code of Node.js Redis Client", Range{Position{0, 52}, Position{1, 20}}},
				Sentence{"I want to jot down some notes on how it works", Range{Position{1, 22}, Position{1, 67}}},
			},
		},
	}
	for _, test := range tests {
		result := parse(test.Text)

		if len(result) != len(test.Expected) {
			t.Errorf("Expected %d sentences, got %d", len(test.Expected), len(result))
			t.Error(result)
			return
		}

		for i, expectedSentence := range test.Expected {
			if !testSentence(t, result[i], expectedSentence) {
				return
			}
		}
	}

}

func testSentence(t *testing.T, actual, expected Sentence) bool {
	if actual.Text != expected.Text {
		t.Errorf("Expected %s, got %s", expected.Text, actual.Text)
		return false
	}

	if !testPosition(t, expected.Text, actual.Range.Start, expected.Range.Start) {
		return false
	}
	if !testPosition(t, expected.Text, actual.Range.End, expected.Range.End) {
		return false
	}
	return true
}

func testPosition(t *testing.T, text string, actual, expected Position) bool {
	if actual.Line != expected.Line {
		t.Errorf("Expected line %s: %d, got %d", text, expected.Line, actual.Line)
		return false
	}
	if actual.Character != expected.Character {
		t.Errorf("Expected character %s: %d, got %d", text, expected.Character, actual.Character)
		return false
	}

	return true
}
