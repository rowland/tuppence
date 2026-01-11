package ast

// Comment represents a comment node
type Comment struct {
	BaseNode
	Text string // Comment text
}

// NewComment creates a new Comment node
func NewComment(text string) *Comment {
	return &Comment{
		BaseNode: BaseNode{Type: NodeComment},
		Text:     text,
	}
}

// String returns a textual representation of the comment
func (c *Comment) String() string {
	return "# " + c.Text
}
