package huffman

import (
	"bytes"
	"encoding/binary"
	"errors"
	"datacompression/pkg/bitstream"
)

// Decompress decodes a Huffman compressed stream back to original bytes.
func Decompress(compressed []byte) ([]byte, error) {
	if len(compressed) < len(MagicHeader)+8 {
		return nil, errors.New("corrupted or truncated archive header")
	}

	if !bytes.Equal(compressed[:len(MagicHeader)], MagicHeader) {
		return nil, errors.New("invalid file signature: not a valid Huffman archive")
	}

	offset := len(MagicHeader)
	r := bytes.NewReader(compressed[offset:])

	var origLen uint64
	if err := binary.Read(r, binary.BigEndian, &origLen); err != nil {
		return nil, err
	}

	br := bitstream.NewBitReader(r)

	root, err := DeserializeTree(br)
	if err != nil {
		return nil, err
	}

	output := make([]byte, 0, origLen)
	curr := root

	for uint64(len(output)) < origLen {
		bit, err := br.ReadBit()
		if err != nil {
			return nil, errors.New("unexpected EOF during decompression")
		}

		if bit == 0 {
			curr = curr.Left
		} else {
			curr = curr.Right
		}

		if curr == nil {
			return nil, errors.New("corrupt tree traversal: reached null node")
		}

		if curr.IsLeaf() {
			output = append(output, curr.Byte)
			curr = root
		}
	}

	return output, nil
}
