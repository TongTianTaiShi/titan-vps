package types

// NodeType node type
type NodeType int

const (
	NodeUnknown NodeType = iota

	NodeTransaction
	NodeBasis
)

func (n NodeType) String() string {
	switch n {
	case NodeTransaction:
		return "transaction"
	case NodeBasis:
		return "basis"
	}

	return ""
}

// RunningNodeType represents the type of the running node.
var RunningNodeType NodeType
