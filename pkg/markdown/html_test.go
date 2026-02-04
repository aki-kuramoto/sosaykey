package markdown

import (
	"testing"
)

func TestDocument_ToHTML(t *testing.T) {
	// Construct an AST manually to test rendering isolated from parsing
	doc := &Document{
		Children: []Node{
			&Header{
				Level: 1,
				Children: []Node{
					&Text{Content: "Title"},
				},
			},
			&Paragraph{
				Children: []Node{
					&Text{Content: "Hello "},
					&Strong{
						Children: []Node{
							&Text{Content: "World"},
						},
					},
					&Text{Content: "!"},
				},
			},
			&List{
				Ordered: false,
				Children: []Node{
					&ListItem{
						Children: []Node{
							&Text{Content: "Item 1"},
						},
					},
					&ListItem{
						Children: []Node{
							&Text{Content: "Item 2"},
						},
					},
				},
			},
			&Link{
				URL: "http://example.com",
				Children: []Node{
					&Text{Content: "Example"},
				},
			},
		},
	}

	html := doc.ToHTML()
	expected := "<h1>Title</h1>\n<p>Hello <strong>World</strong>!</p>\n<ul>\n<li>Item 1</li>\n<li>Item 2</li>\n</ul>\n<a href=\"http://example.com\">Example</a>"

	if html != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, html)
	}
}
