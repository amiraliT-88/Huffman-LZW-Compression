package rle

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var MagicHeader = []byte{'D', 'C', 0x03}

func Compress(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("empty input")
	}

	var buf bytes.Buffer
	buf.Write(MagicHeader)
	origLen := uint64(len(input))
	binary.Write(&buf, binary.BigEndian, origLen)

	n := len(input)
	i := 0
	for i < n {
		b := input[i]
		count := 1
		for i+count < n && input[i+count] == b && count < 255 {
			count++
		}
		buf.WriteByte(b)
		buf.WriteByte(byte(count))
		i += count
	}

	return buf.Bytes(), nil
}

func Decompress(compressed []byte) ([]byte, error) {
	if len(compressed) < len(MagicHeader)+8 {
		return nil, errors.New("corrupted or truncated archive")
	}

	if !bytes.Equal(compressed[:len(MagicHeader)], MagicHeader) {
		return nil, errors.New("invalid rle signature")
	}

	offset := len(MagicHeader)
	r := bytes.NewReader(compressed[offset:])

	var origLen uint64
	if err := binary.Read(r, binary.BigEndian, &origLen); err != nil {
		return nil, err
	}

	payload := compressed[offset+8:]
	output := make([]byte, 0, origLen)

	for i := 0; i < len(payload); i += 2 {
		if i+1 >= len(payload) {
			break
		}
		b := payload[i]
		count := int(payload[i+1])
		for c := 0; c < count; c++ {
			output = append(output, b)
		}
	}

	return output, nil
}
