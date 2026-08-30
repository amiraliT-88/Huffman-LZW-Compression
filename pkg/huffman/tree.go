package huffman

import (
	"container/heap"
	"datacompression/pkg/bitstream"
)

// BuildTree constructs a canonical Huffman tree from a frequency table.
func BuildTree(freqs map[byte]int) *Node {
	if len(freqs) == 0 {
		return nil
	}

	pq := make(PriorityQueue, 0, len(freqs))
	for b, count := range freqs {
		if count > 0 {
			pq = append(pq, &Node{Byte: b, Freq: count})
		}
	}

	if len(pq) == 0 {
		return nil
	}

	// Single unique character edge-case: wrap in symmetric binary structure
	if len(pq) == 1 {
		leaf := pq[0]
		return &Node{
			Freq:  leaf.Freq,
			Left:  &Node{Byte: leaf.Byte, Freq: leaf.Freq},
			Right: &Node{Byte: 0, Freq: 0},
		}
	}

	heap.Init(&pq)

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)

		parent := &Node{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(&pq, parent)
	}

	return heap.Pop(&pq).(*Node)
}

// BuildCodeTable generates prefix binary codes for all leaf symbols in the tree.
func BuildCodeTable(root *Node) map[byte]string {
	table := make(map[byte]string)
	var traverse func(node *Node, path string)
	traverse = func(node *Node, path string) {
		if node == nil {
			return
		}
		if node.IsLeaf() {
			if path == "" {
				path = "0"
			}
			table[node.Byte] = path
			return
		}
		traverse(node.Left, path+"0")
		traverse(node.Right, path+"1")
	}
	traverse(root, "")
	return table
}

// SerializeTree serializes the tree topology using pre-order traversal:
// 0 for internal node, 1 followed by 8-bit byte for leaf node.
func SerializeTree(node *Node, bw *bitstream.BitWriter) error {
	if node == nil {
		return nil
	}
	if node.IsLeaf() {
		if err := bw.WriteBit(1); err != nil {
			return err
		}
		return bw.WriteByte(node.Byte)
	}
	if err := bw.WriteBit(0); err != nil {
		return err
	}
	if err := SerializeTree(node.Left, bw); err != nil {
		return err
	}
	return SerializeTree(node.Right, bw)
}

// DeserializeTree reconstructs the tree topology from the bitstream.
func DeserializeTree(br *bitstream.BitReader) (*Node, error) {
	bit, err := br.ReadBit()
	if err != nil {
		return nil, err
	}

	if bit == 1 {
		val, err := br.ReadByte()
		if err != nil {
			return nil, err
		}
		return &Node{Byte: val}, nil
	}

	left, err := DeserializeTree(br)
	if err != nil {
		return nil, err
	}
	right, err := DeserializeTree(br)
	if err != nil {
		return nil, err
	}

	return &Node{Left: left, Right: right}, nil
}
