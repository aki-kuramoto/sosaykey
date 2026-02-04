package markdown

import (
	"fmt"
	"strings"
)

// ToHTML converts the AST to an HTML string.
func (d *Document) ToHTML() string {
	var sb strings.Builder
	for _, child := range d.Children {
		sb.WriteString(renderNode(child))
	}
	return sb.String()
}

func renderNode(node Node) string {
	switch n := node.(type) {
	case *Header:
		content := renderChildren(n.Children)
		return fmt.Sprintf("<h%d>%s</h%d>\n", n.Level, content, n.Level)
	case *Paragraph:
		content := renderChildren(n.Children)
		return fmt.Sprintf("<p>%s</p>\n", content)
	case *List:
		tag := "ul"
		if n.Ordered {
			tag = "ol"
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("<%s>\n", tag))
		for _, child := range n.Children {
			sb.WriteString(renderNode(child))
		}
		sb.WriteString(fmt.Sprintf("</%s>\n", tag))
		return sb.String()
	case *ListItem:
		content := renderChildren(n.Children)
		return fmt.Sprintf("<li>%s</li>\n", content)
	case *Strong:
		content := renderChildren(n.Children)
		return fmt.Sprintf("<strong>%s</strong>", content)
	case *Emphasis:
		content := renderChildren(n.Children)
		return fmt.Sprintf("<em>%s</em>", content)
	case *Link:
		content := renderChildren(n.Children)
		return fmt.Sprintf("<a href=\"%s\">%s</a>", n.URL, content)
	case *Image:
		return fmt.Sprintf("<img src=\"%s\" alt=\"%s\">", n.URL, n.AltText)
	case *Text:
		return escapeHTML(n.Content)
	default:
		return ""
	}
}

func renderChildren(children []Node) string {
	var sb strings.Builder
	for _, child := range children {
		sb.WriteString(renderNode(child))
	}
	return sb.String()
}

func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&#39;")
	return text
}
