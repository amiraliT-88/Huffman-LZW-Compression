package analyzer

import (
	"fmt"
	"datacompression/pkg/huffman"
	"datacompression/pkg/lzw"
	"datacompression/pkg/rle"
	"time"
)

type BenchResult struct {
	Algorithm       string
	OriginalSize    int
	CompressedSize  int
	Ratio           float64
	CompressionTime time.Duration
	SpeedMBs        float64
}

func RunBenchmark(data []byte) []BenchResult {
	algos := []struct {
		name string
		fn   func([]byte) ([]byte, error)
	}{
		{"Huffman Coding", huffman.Compress},
		{"LZW Dictionary", lzw.Compress},
		{"Run-Length (RLE)", rle.Compress},
	}

	results := make([]BenchResult, 0, len(algos))
	origLen := len(data)
	sizeMB := float64(origLen) / (1024 * 1024)

	for _, a := range algos {
		start := time.Now()
		comp, err := a.fn(data)
		elapsed := time.Since(start)

		if err != nil {
			continue
		}

		compLen := len(comp)
		ratio := (1.0 - (float64(compLen) / float64(origLen))) * 100.0
		speed := sizeMB / elapsed.Seconds()
		if elapsed.Seconds() == 0 {
			speed = sizeMB / 0.000001
		}

		results = append(results, BenchResult{
			Algorithm:       a.name,
			OriginalSize:    origLen,
			CompressedSize:  compLen,
			Ratio:           ratio,
			CompressionTime: elapsed,
			SpeedMBs:        speed,
		})
	}

	return results
}

func FormatBenchmarkTable(results []BenchResult) string {
	header := fmt.Sprintf("\n%-18s | %-12s | %-12s | %-10s | %-12s | %-12s\n",
		"Algorithm", "Original", "Compressed", "Space Saved", "Elapsed", "Throughput")
	sep := "-------------------+--------------+--------------+------------+--------------+--------------\n"
	rows := ""
	for _, r := range results {
		rows += fmt.Sprintf("%-18s | %10d B | %10d B | %9.2f%% | %12s | %9.2f MB/s\n",
			r.Algorithm, r.OriginalSize, r.CompressedSize, r.Ratio, r.CompressionTime, r.SpeedMBs)
	}
	return header + sep + rows
}
