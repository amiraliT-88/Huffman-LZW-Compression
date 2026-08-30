package huffman

type Node struct {
	Byte  byte
	Freq  int
	Left  *Node
	Right *Node
}

func (n *Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}
