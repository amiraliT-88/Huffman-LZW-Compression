package bitstream

import (
	"io"
)

type BitWriter struct {
	w      io.Writer
	buffer byte
	bitPos uint8
}

func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{w: w}
}

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

func (bw *BitWriter) WriteBits(val uint64, nBits int) error {
	for i := nBits - 1; i >= 0; i-- {
		bit := byte((val >> i) & 1)
		if err := bw.WriteBit(bit); err != nil {
			return err
		}
	}
	return nil
}

func (bw *BitWriter) WriteByte(b byte) error {
	return bw.WriteBits(uint64(b), 8)
}

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
