package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"datacompression/pkg/analyzer"
	"datacompression/pkg/huffman"
	"datacompression/pkg/lzw"
	"datacompression/pkg/rle"
)

func hashData(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func main() {
	compress := flag.Bool("c", false, "compress mode")
	decompress := flag.Bool("d", false, "decompress mode")
	analyze := flag.Bool("analyze", false, "analyze Shannon information entropy")
	benchmark := flag.Bool("benchmark", false, "benchmark all compression algorithms")

	inputPath := flag.String("in", "", "input file path")
	outputPath := flag.String("out", "", "output file path")
	algo := flag.String("algo", "huffman", "algorithm: huffman, lzw, rle")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] -in=<file>\n\nFlags:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if *inputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", *inputPath, err)
		os.Exit(1)
	}

	if *analyze {
		res := analyzer.AnalyzeData(data)
		fmt.Println(res.FormatReport())
		return
	}

	if *benchmark {
		fmt.Printf("benchmarking %s (%d bytes)...\n", filepath.Base(*inputPath), len(data))
		results := analyzer.RunBenchmark(data)
		fmt.Println(analyzer.FormatBenchmarkTable(results))
		return
	}

	if *compress {
		outPath := *outputPath
		if outPath == "" {
			outPath = *inputPath + ".dca"
		}

		selected := strings.ToLower(*algo)
		start := time.Now()
		var compressed []byte
		var compErr error

		switch selected {
		case "huffman":
			compressed, compErr = huffman.Compress(data)
		case "lzw":
			compressed, compErr = lzw.Compress(data)
		case "rle":
			compressed, compErr = rle.Compress(data)
		default:
			fmt.Fprintf(os.Stderr, "unknown algorithm %q (use: huffman, lzw, rle)\n", *algo)
			os.Exit(1)
		}

		if compErr != nil {
			fmt.Fprintf(os.Stderr, "compression failed: %v\n", compErr)
			os.Exit(1)
		}
		elapsed := time.Since(start)

		if err := os.WriteFile(outPath, compressed, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", outPath, err)
			os.Exit(1)
		}

		ratio := (1.0 - (float64(len(compressed)) / float64(len(data)))) * 100.0
		fmt.Printf("compressed %s -> %s\n", filepath.Base(*inputPath), filepath.Base(outPath))
		fmt.Printf("algorithm: %s\n", selected)
		fmt.Printf("size:      %d B -> %d B (%.2f%% saved)\n", len(data), len(compressed), ratio)
		fmt.Printf("checksum:  sha256:%s\n", hashData(data)[:16])
		fmt.Printf("elapsed:   %v\n", elapsed)
		return
	}

	if *decompress {
		outPath := *outputPath
		if outPath == "" {
			if strings.HasSuffix(*inputPath, ".dca") {
				outPath = strings.TrimSuffix(*inputPath, ".dca")
			} else {
				outPath = *inputPath + ".extracted"
			}
		}

		start := time.Now()
		var decompressed []byte
		var decompErr error

		if len(data) >= 3 && data[0] == 'D' && data[1] == 'C' {
			switch data[2] {
			case 0x01:
				decompressed, decompErr = huffman.Decompress(data)
			case 0x02:
				decompressed, decompErr = lzw.Decompress(data)
			case 0x03:
				decompressed, decompErr = rle.Decompress(data)
			default:
				decompErr = fmt.Errorf("unsupported format version byte: 0x%02X", data[2])
			}
		} else {
			decompErr = fmt.Errorf("invalid header: not a valid DCA file")
		}

		if decompErr != nil {
			fmt.Fprintf(os.Stderr, "decompression failed: %v\n", decompErr)
			os.Exit(1)
		}
		elapsed := time.Since(start)

		if err := os.WriteFile(outPath, decompressed, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", outPath, err)
			os.Exit(1)
		}

		fmt.Printf("decompressed %s -> %s\n", filepath.Base(*inputPath), filepath.Base(outPath))
		fmt.Printf("restored:   %d bytes\n", len(decompressed))
		fmt.Printf("checksum:   sha256:%s\n", hashData(decompressed)[:16])
		fmt.Printf("elapsed:    %v\n", elapsed)
		return
	}

	flag.Usage()
}
