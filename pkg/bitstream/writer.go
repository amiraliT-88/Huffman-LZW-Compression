package bitstream

import (
	"io"
)

// BitWriter writes individual bits or bit-fields to an underlying io.Writer.
type BitWriter struct {
	w      io.Writer
	buffer byte
	bitPos uint8 // 0 to 7
}

// NewBitWriter initializes a new BitWriter.
func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{w: w}
}

// WriteBit writes a single bit (0 or 1) to the buffer.
func (bw *BitWriter) WriteBit(bit byte) error {
	if bit > 0 {
		bw.buffer |= (1 << (7 - bw.bitPos))
	}
	bw.bitPos++
	if bw.bitPos == 8 {
		if _, err := bw.w.Write([]byte{bw.buffer}); err != nil {
			return err
		}
		bw.buffer = 0
		bw.bitPos = 0
	}
	return nil
}

// WriteBits writes the n least significant bits of val (most-significant bit first).
func (bw *BitWriter) WriteBits(val uint64, nBits int) error {
	for i := nBits - 1; i >= 0; i-- {
		bit := byte((val >> i) & 1)
		if err := bw.WriteBit(bit); err != nil {
			return err
		}
	}
	return nil
}

// WriteByte writes a full 8-bit byte.
func (bw *BitWriter) WriteByte(b byte) error {
	return bw.WriteBits(uint64(b), 8)
}

// Flush writes any remaining buffered bits, padding the final byte with zeros.
// Returns the number of padding bits added (0 to 7).
func (bw *BitWriter) Flush() (uint8, error) {
	if bw.bitPos == 0 {
		return 0, nil
	}
	padding := 8 - bw.bitPos
	if _, err := bw.w.Write([]byte{bw.buffer}); err != nil {
		return 0, err
	}
	bw.buffer = 0
	bw.bitPos = 0
	return padding, nil
}
