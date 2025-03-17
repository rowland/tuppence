package ast

// Comment is the base type for all comment nodes
type Comment struct {
	BaseNode
	Text string // Comment text
}

// NewComment creates a new Comment node
func NewComment(text string) *Comment {
	return &Comment{
		BaseNode: BaseNode{NodeType: NodeComment},
		Text:     text,
	}
}

// String returns a textual representation of the comment
func (c *Comment) String() string {
	return "// " + c.Text
}

// DocComment represents a documentation comment
type DocComment struct {
	BaseNode
	Text string // Comment text
}

// NewDocComment creates a new DocComment node
func NewDocComment(text string) *DocComment {
	return &DocComment{
		BaseNode: BaseNode{NodeType: NodeDocComment},
		Text:     text,
	}
}

// String returns a textual representation of the documentation comment
func (c *DocComment) String() string {
	return "/// " + c.Text
}

// LineComment represents a single-line comment
type LineComment struct {
	BaseNode
	Text string // Comment text
}

// NewLineComment creates a new LineComment node
func NewLineComment(text string) *LineComment {
	return &LineComment{
		BaseNode: BaseNode{NodeType: NodeLineComment},
		Text:     text,
	}
}

// String returns a textual representation of the line comment
func (c *LineComment) String() string {
	return "// " + c.Text
}

// BlockComment represents a multi-line block comment
type BlockComment struct {
	BaseNode
	Text string // Comment text
}

// NewBlockComment creates a new BlockComment node
func NewBlockComment(text string) *BlockComment {
	return &BlockComment{
		BaseNode: BaseNode{NodeType: NodeBlockComment},
		Text:     text,
	}
}

// String returns a textual representation of the block comment
func (c *BlockComment) String() string {
	return "/* " + c.Text + " */"
}
