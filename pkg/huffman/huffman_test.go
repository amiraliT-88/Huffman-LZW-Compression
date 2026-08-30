package huffman

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestHuffmanRoundtrip(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "Classic Text",
			data: []byte("Data compression is the science of encoding information using fewer bits than original representation."),
		},
		{
			name: "Repetitive Stream",
			data: bytes.Repeat([]byte("HuffmanCodingAlgorithmInGoWithZeroDependencies"), 80),
		},
		{
			name: "All 256 Byte Symbols",
			data: func() []byte {
				b := make([]byte, 256)
				for i := 0; i < 256; i++ {
					b[i] = byte(i)
				}
				return b
			}(),
		},
		{
			name: "Single Character Repeated",
			data: bytes.Repeat([]byte("A"), 100),
		},
		{
			name: "Random Binary Block",
			data: func() []byte {
				buf := make([]byte, 2048)
				rand.Read(buf)
				return buf
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			compressed, err := Compress(tc.data)
			if err != nil {
				t.Fatalf("Compression failed: %v", err)
			}

			decompressed, err := Decompress(compressed)
			if err != nil {
				t.Fatalf("Decompression failed: %v", err)
			}

			if !bytes.Equal(decompressed, tc.data) {
				t.Fatalf("Decompressed data mismatch! Original len=%d, Restored len=%d", len(tc.data), len(decompressed))
			}
		})
	}
}
