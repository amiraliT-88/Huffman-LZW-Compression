package huffman

import (
	"bytes"
	"encoding/binary"
	"errors"
	"datacompression/pkg/bitstream"
)

var (
	MagicHeader = []byte{'D', 'C', 0x01} // Data-Compression Huffman format
)

// Compress encodes data using Huffman variable-length prefix coding.
func Compress(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("cannot compress empty input")
	}

	freqs := make(map[byte]int)
	for _, b := range input {
		freqs[b]++
	}

	tree := BuildTree(freqs)
	if tree == nil {
		return nil, errors.New("failed to construct Huffman tree")
	}
	codeTable := BuildCodeTable(tree)

	var buf bytes.Buffer

	// Write Magic Header
	if _, err := buf.Write(MagicHeader); err != nil {
		return nil, err
	}

	// Write Original Length (uint64, 8 bytes)
	origLen := uint64(len(input))
	if err := binary.Write(&buf, binary.BigEndian, origLen); err != nil {
		return nil, err
	}

	bw := bitstream.NewBitWriter(&buf)

	// Serialize Tree topology in bitstream
	if err := SerializeTree(tree, bw); err != nil {
		return nil, err
	}

	// Encode payload bits
	for _, b := range input {
		code := codeTable[b]
		for i := 0; i < len(code); i++ {
			bit := byte(code[i] - '0')
			if err := bw.WriteBit(bit); err != nil {
				return nil, err
			}
		}
	}

	// Flush remaining bits
	if _, err := bw.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
