package huffman

import (
	"bytes"
	"encoding/binary"
	"errors"

	"datacompression/pkg/bitstream"
)

var MagicHeader = []byte{'D', 'C', 0x01}

func Compress(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("empty input")
	}

	freqs := make(map[byte]int)
	for _, b := range input {
		freqs[b]++
	}

	tree := BuildTree(freqs)
	if tree == nil {
		return nil, errors.New("failed to build tree")
	}
	codeTable := BuildCodeTable(tree)

	var buf bytes.Buffer

	if _, err := buf.Write(MagicHeader); err != nil {
		return nil, err
	}

	origLen := uint64(len(input))
	if err := binary.Write(&buf, binary.BigEndian, origLen); err != nil {
		return nil, err
	}

	bw := bitstream.NewBitWriter(&buf)
	if err := SerializeTree(tree, bw); err != nil {
		return nil, err
	}

	for _, b := range input {
		code := codeTable[b]
		for i := 0; i < len(code); i++ {
			bit := byte(code[i] - '0')
			if err := bw.WriteBit(bit); err != nil {
				return nil, err
			}
		}
	}

	if _, err := bw.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
