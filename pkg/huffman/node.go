package huffman

// Node represents a vertex in the Huffman binary prefix tree.
type Node struct {
	Byte  byte
	Freq  int
	Left  *Node
	Right *Node
}

// IsLeaf returns true if the node represents a terminal symbol.
func (n *Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}
