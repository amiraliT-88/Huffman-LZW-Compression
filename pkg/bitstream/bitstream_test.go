package bitstream

import (
	"bytes"
	"testing"
)

func TestBitStreamRoundtrip(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)

	bits := []byte{1, 0, 1, 1, 0}
	for _, b := range bits {
		if err := bw.WriteBit(b); err != nil {
			t.Fatalf("WriteBit failed: %v", err)
		}
	}

	val := uint64(0xABC)
	if err := bw.WriteBits(val, 12); err != nil {
		t.Fatalf("WriteBits failed: %v", err)
	}

	padding, err := bw.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}
	if padding != 7 {
		t.Errorf("Expected 7 padding bits, got %d", padding)
	}

	br := NewBitReader(&buf)
	for i, expected := range bits {
		bit, err := br.ReadBit()
		if err != nil {
			t.Fatalf("ReadBit %d failed: %v", i, err)
		}
		if bit != expected {
			t.Errorf("Bit %d: expected %d, got %d", i, expected, bit)
		}
	}

	readVal, err := br.ReadBits(12)
	if err != nil {
		t.Fatalf("ReadBits failed: %v", err)
	}
	if readVal != val {
		t.Errorf("Expected 0x%X, got 0x%X", val, readVal)
	}
}
