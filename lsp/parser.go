package lsp

import (
	"strings"
)

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Sentence struct {
	Text  string `json:"text"`
	Range Range  `json:"range"`
}

func (s *Server) parse(text string) []Sentence {
	lines := strings.Split(text, "\n")

	remaining := Sentence{}
	result := []Sentence{}

	frontMatterChecked := false
	skip := false

	for lineNumber, line := range lines {
		// Ignore front matter
		if line == "+++" || line == "---" {
			if lineNumber == 0 {
				frontMatterChecked = !frontMatterChecked
			}
			skip = !skip
			continue
		}

		if !frontMatterChecked {
			continue
		}

		// Ignore code blocks
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			skip = !skip
			continue
		}

		// Ignore links and HTML comments
		if strings.Contains(line, "]: http") || strings.HasPrefix(line, "<!-") || skip {
			continue
		}

		if strings.Trim(line, " ") == "" {
			// if len(paragraph) > 0 {
			// 	paragraphs = append(paragraphs, paragraph)
			// 	paragraph = ""
			// }
			continue
		}

		character := 0
		parts := paragraphRegex.Split(line, -1)
		for i, part := range parts {
			start := Position{Line: lineNumber, Character: character}
			end := Position{Line: lineNumber, Character: character + len(part)}
			if remaining.Text != "" {
				remaining.Range.End.Character = len(part)
				remaining.Range.End.Line = lineNumber
				remaining.Text += " " + part

				result = append(result, remaining)
				remaining = Sentence{}
			}
			if i+1 == len(parts) {
				remaining = Sentence{Text: part, Range: Range{Start: start, End: end}}
			} else {
				result = append(result, Sentence{
					Text:  part,
					Range: Range{Start: start, End: end},
				})
				character += len(part) + 2
			}
		}
	}
	return result
}
