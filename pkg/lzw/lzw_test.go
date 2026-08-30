package lzw

import (
	"bytes"
	"testing"
)

func TestLZWRoundtrip(t *testing.T) {
	data := []byte("TOBEORNOTTOBEORTOBEORNOTTOBE_LZW_ALGORITHM_TEST_DATA_IN_GO")
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
