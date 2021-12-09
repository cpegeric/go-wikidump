package wikitext

import (
	"regexp"
	"strings"
)

// We are going to iterate over the characters in text one by one and check multiple condition as to where
// the characters need to be written
// TODO: remove blockquote
func ToPlain(text string) string {
	recTemplate, recLink := false, false
	curBrCount, sqBrCount := 0, 0
	plainIndex := 0
	var linkBuilder, plainBuilder strings.Builder
	for _, char := range text {
		switch char {
		case '\'':
			continue
		case '{':
			if !recTemplate {
				recTemplate = true
			} else {
				curBrCount++
			}
		case '}':
			if curBrCount == 0 {
				recTemplate = false
			} else {
				curBrCount--
			}
		case '[':
			if recTemplate {
			} else if !recLink {
				recLink = true
				linkBuilder.WriteRune('[')
			} else {
				linkBuilder.WriteRune('[')
				sqBrCount++
			}
		case ']':
			if recTemplate {
			} else if sqBrCount == 0 {
				recLink = false
				linkBuilder.WriteRune(']')
				linkString := linkBuilder.String()
				plainBuilder.WriteString(linkDisplay(linkString))
				plainIndex += len(linkString)
				linkBuilder = strings.Builder{}
			} else {
				linkBuilder.WriteRune(char)
				sqBrCount--
			}
		default:
			if recLink {
				linkBuilder.WriteRune(char)
			}
			if !recLink && !recTemplate {
				plainBuilder.WriteRune(char)
				plainIndex++
			}
		}
	}
	return plainBuilder.String()
}

// TODO: handle nested links
func linkDisplay(link string) string {
	re := regexp.MustCompile(`[\[\]]`)
	return re.ReplaceAllLiteralString(link, "")
}

// TODO: trim empty lines and lines with only * in them.
// TODO: split sections.
// TODO: filtering by sections.
