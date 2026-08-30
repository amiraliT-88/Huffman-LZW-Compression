package lzw

import (
	"bytes"
	"encoding/binary"
	"errors"
	"datacompression/pkg/bitstream"
)

var (
	MagicHeader = []byte{'D', 'C', 0x02} // Data-Compression LZW format
)

const (
	MaxDictSize = 4096
	CodeBits    = 12
)

// Compress encodes data using the LZW dictionary-based algorithm.
func Compress(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("cannot compress empty input")
	}

	var buf bytes.Buffer
	if _, err := buf.Write(MagicHeader); err != nil {
		return nil, err
	}
	origLen := uint64(len(input))
	if err := binary.Write(&buf, binary.BigEndian, origLen); err != nil {
		return nil, err
	}

	bw := bitstream.NewBitWriter(&buf)

	dict := make(map[string]uint16)
	for i := 0; i < 256; i++ {
		dict[string([]byte{byte(i)})] = uint16(i)
	}
	nextCode := uint16(256)

	var current string
	for _, b := range input {
		combined := current + string([]byte{b})
		if _, exists := dict[combined]; exists {
			current = combined
		} else {
			if err := bw.WriteBits(uint64(dict[current]), CodeBits); err != nil {
				return nil, err
			}
			if nextCode < MaxDictSize {
				dict[combined] = nextCode
				nextCode++
			}
			current = string([]byte{b})
		}
	}

	if current != "" {
		if err := bw.WriteBits(uint64(dict[current]), CodeBits); err != nil {
			return nil, err
		}
	}

	if _, err := bw.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
