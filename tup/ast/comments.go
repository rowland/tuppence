package ast

type Comment struct {
	BaseNode
	Text string // Comment text
}

func NewComment(text string) *Comment {
	return &Comment{
		BaseNode: BaseNode{Type: NodeComment},
		Text:     text,
	}
}

func (c *Comment) String() string {
	return "# " + c.Text
}
