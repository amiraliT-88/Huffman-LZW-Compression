package rle

import (
	"bytes"
	"testing"
)

func TestRLERoundtrip(t *testing.T) {
	data := []byte("AAAAABBBBBCCCCCCCDDDDDDEEEEEEEFFFFFF_RLE_DATA")
	compressed, err := Compress(data)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}

	decompressed, err := Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}

	if !bytes.Equal(decompressed, data) {
		t.Fatalf("Mismatch! Expected %s, got %s", string(data), string(decompressed))
	}
}
