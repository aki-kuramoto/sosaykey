package markdown

import (
	"testing"
)

func TestParser_Parse_Header(t *testing.T) {
	input := "# Header 1\n## Header 2\n### Header 3\n#### Header 4\n##### Header 5\n###### Header 6"
	parser := NewParser(input)
	doc := parser.Parse()

	if len(doc.Children) != 6 {
		t.Errorf("Expected 6 children, got %d", len(doc.Children))
	}

	expectedLevels := []int{1, 2, 3, 4, 5, 6}
	for i, node := range doc.Children {
		header, ok := node.(*Header)
		if !ok {
			t.Errorf("Node %d is not a Header, got %T", i, node)
			continue
		}
		if header.Level != expectedLevels[i] {
			t.Errorf("Node %d level expected %d, got %d", i, expectedLevels[i], header.Level)
		}

		expectedContent := "Header " + string(rune('0'+expectedLevels[i]))
		if len(header.Children) != 1 {
			t.Errorf("Header %d expected 1 child, got %d", i, len(header.Children))
		}
		text, ok := header.Children[0].(*Text)
		if !ok {
			t.Errorf("Header %d child is not Text", i)
		}
		if text.Content != expectedContent {
			t.Errorf("Header %d content expected %q, got %q", i, expectedContent, text.Content)
		}
	}
}

func TestParser_Parse_Paragraph(t *testing.T) {
	input := "This is a paragraph.\nAnother paragraph."
	parser := NewParser(input)
	doc := parser.Parse()

	if len(doc.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(doc.Children))
	}

	expectedTexts := []string{"This is a paragraph.", "Another paragraph."}
	for i, node := range doc.Children {
		p, ok := node.(*Paragraph)
		if !ok {
			t.Errorf("Node %d is not a Paragraph, got %T", i, node)
			continue
		}
		if len(p.Children) != 1 {
			t.Errorf("Paragraph %d expected 1 child, got %d", i, len(p.Children))
		}
		text, ok := p.Children[0].(*Text)
		if !ok {
			t.Errorf("Paragraph %d child is not Text", i)
		}
		if text.Content != expectedTexts[i] {
			t.Errorf("Paragraph %d content expected %q, got %q", i, expectedTexts[i], text.Content)
		}
	}
}

func TestParser_Parse_UnorderedList(t *testing.T) {
	input := "- Item 1\n* Item 2\n+ Item 3"
	parser := NewParser(input)
	doc := parser.Parse()

	if len(doc.Children) != 1 {
		t.Fatalf("Expected 1 child (List), got %d", len(doc.Children))
	}

	list, ok := doc.Children[0].(*List)
	if !ok {
		t.Fatalf("Expected List node, got %T", doc.Children[0])
	}

	if list.Ordered {
		t.Error("Expected unordered list")
	}

	if len(list.Children) != 3 {
		t.Fatalf("Expected 3 list items, got %d", len(list.Children))
	}
}

func TestParser_Parse_InlineElements(t *testing.T) {
	input := "Normal **Bold** *Italic* [Link](http://example.com) ![Image](http://example.com/img.png)"
	parser := NewParser(input)
	doc := parser.Parse()

	if len(doc.Children) != 1 {
		t.Fatalf("Expected 1 child (Paragraph), got %d", len(doc.Children))
	}

	p, ok := doc.Children[0].(*Paragraph)
	if !ok {
		t.Fatalf("Expected Paragraph, got %T", doc.Children[0])
	}

	// Expected children:
	// 1. Text "Normal "
	// 2. Strong (Child: Text "Bold")
	// 3. Text " "
	// 4. Emphasis (Child: Text "Italic")
	// 5. Text " "
	// 6. Link (Child: Text "Link", URL: ...)
	// 7. Text " "
	// 8. Image (Alt: "Image", URL: ...)

	if len(p.Children) != 8 {
		t.Errorf("Expected 8 children in paragraph, got %d", len(p.Children))
		for i, c := range p.Children {
			t.Logf("Child %d: %T %v", i, c, c)
		}
	}

	// Verify Structure briefly
	if _, ok := p.Children[1].(*Strong); !ok {
		t.Errorf("Child 1 should be Strong")
	}
	if _, ok := p.Children[3].(*Emphasis); !ok {
		t.Errorf("Child 3 should be Emphasis")
	}
	if link, ok := p.Children[5].(*Link); !ok {
		t.Errorf("Child 5 should be Link")
	} else if link.URL != "http://example.com" {
		t.Errorf("Link URL match failed")
	}
	if img, ok := p.Children[7].(*Image); !ok {
		t.Errorf("Child 7 should be Image")
	} else if img.URL != "http://example.com/img.png" {
		t.Errorf("Image URL match failed")
	}
}

func TestParser_Parse_NestedInline(t *testing.T) {
	input := "**Bold *Italic* inside**"
	parser := NewParser(input)
	doc := parser.Parse()

	p := doc.Children[0].(*Paragraph)
	if len(p.Children) != 1 {
		t.Fatalf("Expected 1 child (Strong), got %d", len(p.Children))
	}

	strong := p.Children[0].(*Strong)
	// Strong Children: Text "Bold ", Emphasis(Text "Italic"), Text " inside"
	if len(strong.Children) != 3 {
		t.Errorf("Expected 3 children in Strong, got %d", len(strong.Children))
	}

	if _, ok := strong.Children[1].(*Emphasis); !ok {
		t.Errorf("Expected Emphasis inside Strong")
	}
}
