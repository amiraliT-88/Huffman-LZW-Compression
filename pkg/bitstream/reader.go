package bitstream

import (
	"io"
)

// BitReader reads individual bits or bit-fields from an underlying io.Reader.
type BitReader struct {
	r      io.Reader
	buffer byte
	bitPos uint8 // 8 means buffer is exhausted
}

// NewBitReader initializes a new BitReader.
func NewBitReader(r io.Reader) *BitReader {
	return &BitReader{r: r, bitPos: 8}
}

// ReadBit reads a single bit (0 or 1).
func (br *BitReader) ReadBit() (byte, error) {
	if br.bitPos == 8 {
		buf := []byte{0}
		if _, err := io.ReadFull(br.r, buf); err != nil {
			return 0, err
		}
		br.buffer = buf[0]
		br.bitPos = 0
	}
	bit := (br.buffer >> (7 - br.bitPos)) & 1
	br.bitPos++
	return bit, nil
}

// ReadBits reads nBits and returns them as a uint64 word.
func (br *BitReader) ReadBits(nBits int) (uint64, error) {
	var val uint64
	for i := 0; i < nBits; i++ {
		bit, err := br.ReadBit()
		if err != nil {
			return 0, err
		}
		val = (val << 1) | uint64(bit)
	}
	return val, nil
}

// ReadByte reads a full 8-bit byte.
func (br *BitReader) ReadByte() (byte, error) {
	val, err := br.ReadBits(8)
	if err != nil {
		return 0, err
	}
	return byte(val), nil
}
