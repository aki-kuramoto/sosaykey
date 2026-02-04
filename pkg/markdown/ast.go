package markdown

// NodeType represents the type of a markdown node.
type NodeType int

const (
	NodeTypeDocument NodeType = iota
	NodeTypeHeader
	NodeTypeParagraph
	NodeTypeList
	NodeTypeListItem
	NodeTypeCodeBlock
	NodeTypeText
	NodeTypeLink
	NodeTypeImage
	NodeTypeStrong
	NodeTypeEmphasis
)

// Node represents a node in the markdown AST.
type Node interface {
	Type() NodeType
	String() string
}

// Ensure implementations satisfy the interface
var _ Node = (*Document)(nil)
var _ Node = (*Header)(nil)
var _ Node = (*Paragraph)(nil)
var _ Node = (*Text)(nil)

// Document is the root node of the AST.
type Document struct {
	Children []Node
}

func (d *Document) Type() NodeType { return NodeTypeDocument }
func (d *Document) String() string { return "Document" }

// Block Node Types

type Header struct {
	Level    int
	Children []Node
}

func (h *Header) Type() NodeType { return NodeTypeHeader }
func (h *Header) String() string { return "Header" }

type Paragraph struct {
	Children []Node
}

func (p *Paragraph) Type() NodeType { return NodeTypeParagraph }
func (p *Paragraph) String() string { return "Paragraph" }

// Inline Node Types

type Text struct {
	Content string
}

func (t *Text) Type() NodeType { return NodeTypeText }
func (t *Text) String() string { return "Text" }

type List struct {
	Ordered  bool
	Children []Node
}

func (l *List) Type() NodeType { return NodeTypeList }
func (l *List) String() string { return "List" }

type ListItem struct {
	Children []Node
}

func (l *ListItem) Type() NodeType { return NodeTypeListItem }
func (l *ListItem) String() string { return "ListItem" }

type Strong struct {
	Children []Node
}

func (s *Strong) Type() NodeType { return NodeTypeStrong }
func (s *Strong) String() string { return "Strong" }

type Emphasis struct {
	Children []Node
}

func (e *Emphasis) Type() NodeType { return NodeTypeEmphasis }
func (e *Emphasis) String() string { return "Emphasis" }

type Link struct {
	Children []Node
	URL      string
}

func (l *Link) Type() NodeType { return NodeTypeLink }
func (l *Link) String() string { return "Link" }

type Image struct {
	AltText string
	URL     string
}

func (i *Image) Type() NodeType { return NodeTypeImage }
func (i *Image) String() string { return "Image" }
