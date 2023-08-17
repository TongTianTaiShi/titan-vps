package types

// NodeType node type
type NodeType int

const (
	NodeUnknown NodeType = iota

	NodeTransaction
	NodeMall
)

func (n NodeType) String() string {
	switch n {
	case NodeTransaction:
		return "transaction"
	case NodeMall:
		return "mall"
	}

	return ""
}

// RunningNodeType represents the type of the running node.
var RunningNodeType NodeType
