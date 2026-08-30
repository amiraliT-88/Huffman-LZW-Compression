package analyzer

import (
	"fmt"
	"math"
)

type AnalysisResult struct {
	TotalBytes            int
	UniqueSymbols         int
	ShannonEntropy        float64
	MaxTheoreticalRatio   float64
	OptimalCompressedSize int
	FrequencyTable        map[byte]int
}

func AnalyzeData(data []byte) AnalysisResult {
	n := len(data)
	if n == 0 {
		return AnalysisResult{}
	}

	freqs := make(map[byte]int)
	for _, b := range data {
		freqs[b]++
	}

	var entropy float64
	for _, count := range freqs {
		p := float64(count) / float64(n)
		entropy -= p * math.Log2(p)
	}

	maxRatio := (1.0 - (entropy / 8.0)) * 100.0
	if maxRatio < 0 {
		maxRatio = 0
	}

	optimalBytes := int(math.Ceil((entropy * float64(n)) / 8.0))

	return AnalysisResult{
		TotalBytes:            n,
		UniqueSymbols:         len(freqs),
		ShannonEntropy:        entropy,
		MaxTheoreticalRatio:   maxRatio,
		OptimalCompressedSize: optimalBytes,
		FrequencyTable:        freqs,
	}
}

func (r AnalysisResult) FormatReport() string {
	return fmt.Sprintf("File Analysis (Shannon Entropy)\n" +
		"Total Size:              %d bytes\n" +
		"Unique Symbols:          %d / 256\n" +
		"Shannon Info H(X):       %.4f bits/byte\n" +
		"Max Theoretical Saving:  %.2f%%\n" +
		"Theoretical Lower Bound: %d bytes",
		r.TotalBytes, r.UniqueSymbols, r.ShannonEntropy, r.MaxTheoreticalRatio, r.OptimalCompressedSize)
}
