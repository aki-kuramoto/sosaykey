package markdown

import (
	"bufio"
	"strings"
)

// Parser holds the state of the markdown parser.
type Parser struct {
	input string
}

// NewParser creates a new Parser instance.
func NewParser(input string) *Parser {
	return &Parser{input: input}
}

// Parse parses the input markdown and returns the root Document node.
func (p *Parser) Parse() *Document {
	doc := &Document{Children: []Node{}}
	scanner := bufio.NewScanner(strings.NewReader(p.input))

	var currentList *List

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			currentList = nil // Reset list context on empty line
			continue
		}

		isHeader := false
		// Header parsing
		if strings.HasPrefix(trimmedLine, "#") {
			// Initialize indentation
			level := 0
			// Count the number of hashes
			for _, char := range trimmedLine {
				if char == '#' {
					level++
				} else {
					break
				}
			}

			// Valid headers range from 1 to 6
			if level >= 1 && level <= 6 {
				// Ensure there's a space after the hashes
				contentStartIndex := level
				if len(trimmedLine) > level {
					if trimmedLine[level] == ' ' {
						contentStartIndex++

						content := strings.TrimSpace(trimmedLine[contentStartIndex:])
						header := &Header{
							Level:    level,
							Children: p.parseInline(content),
						}
						doc.Children = append(doc.Children, header)
						currentList = nil // Headers break lists
						isHeader = true
					}
				}
			}
		}

		if isHeader {
			continue
		}

		// List parsing
		isUnorderedList := strings.HasPrefix(trimmedLine, "- ") || strings.HasPrefix(trimmedLine, "* ") || strings.HasPrefix(trimmedLine, "+ ")

		isOrderedList := false
		orderedContentStart := 0
		if !isUnorderedList && len(trimmedLine) > 0 && trimmedLine[0] >= '0' && trimmedLine[0] <= '9' {
			// Check for '1. ' pattern
			dotIndex := strings.Index(trimmedLine, ".")
			if dotIndex > 0 && dotIndex+1 < len(trimmedLine) && trimmedLine[dotIndex+1] == ' ' {
				// Verify all chars before dot are digits
				allDigits := true
				for i := 0; i < dotIndex; i++ {
					if trimmedLine[i] < '0' || trimmedLine[i] > '9' {
						allDigits = false
						break
					}
				}
				if allDigits {
					isOrderedList = true
					orderedContentStart = dotIndex + 2
				}
			}
		}

		if isUnorderedList || isOrderedList {
			// Determine if we should start a new list or continue existing
			startNewList := false
			if currentList == nil {
				startNewList = true
			} else {
				// Check if list type matches
				if currentList.Ordered != isOrderedList {
					startNewList = true
				}
			}

			if startNewList {
				currentList = &List{Ordered: isOrderedList, Children: []Node{}}
				doc.Children = append(doc.Children, currentList)
			}

			// Extract content
			var content string
			if isUnorderedList {
				content = strings.TrimSpace(trimmedLine[2:])
			} else {
				content = strings.TrimSpace(trimmedLine[orderedContentStart:])
			}

			listItem := &ListItem{
				Children: p.parseInline(content),
			}
			currentList.Children = append(currentList.Children, listItem)
			continue
		}

		currentList = nil // Paragraphs break lists
		// Default to paragraph
		// TODO: Handle multi-line paragraphs
		paragraph := &Paragraph{
			Children: p.parseInline(trimmedLine),
		}
		doc.Children = append(doc.Children, paragraph)
	}

	return doc
}

// parseInline parses inline elements like bold, italic, links, etc.
func (p *Parser) parseInline(text string) []Node {
	var nodes []Node
	i := 0
	length := len(text)

	for i < length {
		// Check for Bold (**text**)
		if i+1 < length && text[i] == '*' && text[i+1] == '*' {
			end := strings.Index(text[i+2:], "**")
			if end != -1 {
				content := text[i+2 : i+2+end]
				nodes = append(nodes, &Strong{Children: p.parseInline(content)})
				i = i + 2 + end + 2
				continue
			}
		}

		// Check for Italic (*text*)
		if text[i] == '*' {
			// Find closing *
			// Note: This is a very simplified check and doesn't handle all edge cases or escapes
			remaining := text[i+1:]
			end := strings.Index(remaining, "*")
			if end != -1 {
				content := remaining[:end]
				nodes = append(nodes, &Emphasis{Children: p.parseInline(content)})
				i = i + 1 + end + 1
				continue
			}
		}

		// Check for Image (![alt](url))
		if i+1 < length && text[i] == '!' && text[i+1] == '[' {
			closeBracket := strings.Index(text[i:], "]")
			if closeBracket != -1 {
				openParen := strings.Index(text[i+closeBracket:], "(")
				if openParen != -1 && openParen == 1 { // must be directly after ]
					closeParen := strings.Index(text[i+closeBracket+openParen:], ")")
					if closeParen != -1 {
						altText := text[i+2 : i+closeBracket]
						url := text[i+closeBracket+openParen+1 : i+closeBracket+openParen+closeParen]
						nodes = append(nodes, &Image{AltText: altText, URL: url})
						i = i + closeBracket + openParen + closeParen + 1
						continue
					}
				}
			}
		}

		// Check for Link ([text](url))
		if text[i] == '[' {
			closeBracket := strings.Index(text[i:], "]")
			if closeBracket != -1 {
				openParen := strings.Index(text[i+closeBracket:], "(")
				if openParen != -1 && openParen == 1 { // must be directly after ]
					closeParen := strings.Index(text[i+closeBracket+openParen:], ")")
					if closeParen != -1 {
						linkText := text[i+1 : i+closeBracket]
						url := text[i+closeBracket+openParen+1 : i+closeBracket+openParen+closeParen]
						nodes = append(nodes, &Link{Children: p.parseInline(linkText), URL: url})
						i = i + closeBracket + openParen + closeParen + 1
						continue
					}
				}
			}
		}

		// Text
		// Accumulate text until next special char
		start := i
		for i < length {
			if strings.HasPrefix(text[i:], "**") || text[i] == '*' ||
				strings.HasPrefix(text[i:], "![") || text[i] == '[' {
				break
			}
			i++
		}
		if i > start {
			nodes = append(nodes, &Text{Content: text[start:i]})
		}
	}

	return nodes
}
