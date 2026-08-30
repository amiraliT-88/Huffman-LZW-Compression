package lzw

import (
	"bytes"
	"encoding/binary"
	"errors"
	"datacompression/pkg/bitstream"
)

// Decompress decodes an LZW compressed stream back to original bytes.
func Decompress(compressed []byte) ([]byte, error) {
	if len(compressed) < len(MagicHeader)+8 {
		return nil, errors.New("corrupted or truncated archive header")
	}

	if !bytes.Equal(compressed[:len(MagicHeader)], MagicHeader) {
		return nil, errors.New("invalid file signature: not a valid LZW archive")
	}

	offset := len(MagicHeader)
	r := bytes.NewReader(compressed[offset:])

	var origLen uint64
	if err := binary.Read(r, binary.BigEndian, &origLen); err != nil {
		return nil, err
	}

	br := bitstream.NewBitReader(r)

	dict := make(map[uint16]string)
	for i := 0; i < 256; i++ {
		dict[uint16(i)] = string([]byte{byte(i)})
	}
	nextCode := uint16(256)

	firstCode, err := br.ReadBits(CodeBits)
	if err != nil {
		return nil, err
	}

	entry, ok := dict[uint16(firstCode)]
	if !ok {
		return nil, errors.New("corrupt initial LZW code")
	}

	var output bytes.Buffer
	output.WriteString(entry)
	prev := entry

	for uint64(output.Len()) < origLen {
		codeVal, err := br.ReadBits(CodeBits)
		if err != nil {
			break
		}
		code := uint16(codeVal)

		var current string
		if val, exists := dict[code]; exists {
			current = val
		} else if code == nextCode {
			current = prev + string(prev[0])
		} else {
			return nil, errors.New("corrupt LZW codeword sequence")
		}

		output.WriteString(current)

		if nextCode < MaxDictSize {
			dict[nextCode] = prev + string(current[0])
			nextCode++
		}
		prev = current
	}

	return output.Bytes(), nil
}
