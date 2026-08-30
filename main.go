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

const Version = "1.0.0"

func calculateSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func printBanner() {
	fmt.Println("==================================================================")
	fmt.Printf("   Data Compression Engine v%s - Zero-Dependency CS Suite\n", Version)
	fmt.Println("==================================================================")
}

func main() {
	compressCmd := flag.Bool("c", false, "Compress mode")
	decompressCmd := flag.Bool("d", false, "Decompress mode")
	analyzeCmd := flag.Bool("analyze", false, "Analyze Shannon Information Entropy of file")
	benchmarkCmd := flag.Bool("benchmark", false, "Benchmark all algorithms on input file")

	inputPath := flag.String("in", "", "Input file path")
	outputPath := flag.String("out", "", "Output file path (optional for analyze/benchmark)")
	algo := flag.String("algo", "huffman", "Compression algorithm: huffman, lzw, rle")

	flag.Usage = func() {
		printBanner()
		fmt.Println("Usage:")
		fmt.Println("  datacompression -c -algo=[huffman|lzw|rle] -in=<file> -out=<file.dca>")
		fmt.Println("  datacompression -d -in=<file.dca> -out=<file>")
		fmt.Println("  datacompression -analyze -in=<file>")
		fmt.Println("  datacompression -benchmark -in=<file>")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *inputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Printf("[Error] Failed to read input file: %v\n", err)
		os.Exit(1)
	}

	if *analyzeCmd {
		printBanner()
		res := analyzer.AnalyzeData(data)
		fmt.Println(res.FormatReport())
		return
	}

	if *benchmarkCmd {
		printBanner()
		fmt.Printf("Benchmarking %s (%d bytes)...\n", filepath.Base(*inputPath), len(data))
		results := analyzer.RunBenchmark(data)
		fmt.Println(analyzer.FormatBenchmarkTable(results))
		return
	}

	if *compressCmd {
		printBanner()
		outPath := *outputPath
		if outPath == "" {
			outPath = *inputPath + ".dca"
		}

		selectedAlgo := strings.ToLower(*algo)
		fmt.Printf("=> Compressing '%s' using algorithm: %s\n", filepath.Base(*inputPath), strings.ToUpper(selectedAlgo))
		shaIn := calculateSHA256(data)
		fmt.Printf("   Original Size: %d bytes (SHA256: %s...)\n", len(data), shaIn[:12])

		start := time.Now()
		var compressed []byte
		var compErr error

		switch selectedAlgo {
		case "huffman":
			compressed, compErr = huffman.Compress(data)
		case "lzw":
			compressed, compErr = lzw.Compress(data)
		case "rle":
			compressed, compErr = rle.Compress(data)
		default:
			fmt.Printf("[Error] Unknown algorithm '%s'. Choose: huffman, lzw, rle\n", *algo)
			os.Exit(1)
		}

		if compErr != nil {
			fmt.Printf("[Error] Compression failed: %v\n", compErr)
			os.Exit(1)
		}
		elapsed := time.Since(start)

		if err := os.WriteFile(outPath, compressed, 0644); err != nil {
			fmt.Printf("[Error] Failed to save compressed file: %v\n", err)
			os.Exit(1)
		}

		ratio := (1.0 - (float64(len(compressed)) / float64(len(data)))) * 100.0
		fmt.Printf("[SUCCESS] Saved to '%s'\n", outPath)
		fmt.Printf("   Compressed Size: %d bytes\n", len(compressed))
		fmt.Printf("   Space Saved:     %.2f%%\n", ratio)
		fmt.Printf("   Time Elapsed:    %v\n", elapsed)
		return
	}

	if *decompressCmd {
		printBanner()
		outPath := *outputPath
		if outPath == "" {
			if strings.HasSuffix(*inputPath, ".dca") {
				outPath = strings.TrimSuffix(*inputPath, ".dca")
			} else {
				outPath = *inputPath + ".extracted"
			}
		}

		fmt.Printf("=> Decompressing '%s'...\n", filepath.Base(*inputPath))
		start := time.Now()

		var decompressed []byte
		var decompErr error

		if len(data) >= 3 && data[0] == 'D' && data[1] == 'C' {
			switch data[2] {
			case 0x01:
				fmt.Println("   Auto-detected format: Huffman Coding")
				decompressed, decompErr = huffman.Decompress(data)
			case 0x02:
				fmt.Println("   Auto-detected format: LZW Dictionary")
				decompressed, decompErr = lzw.Decompress(data)
			case 0x03:
				fmt.Println("   Auto-detected format: Run-Length (RLE)")
				decompressed, decompErr = rle.Decompress(data)
			default:
				decompErr = fmt.Errorf("unsupported format version byte: 0x%02X", data[2])
			}
		} else {
			decompErr = fmt.Errorf("invalid header: file is not a valid Data-Compression archive")
		}

		if decompErr != nil {
			fmt.Printf("[Error] Decompression failed: %v\n", decompErr)
			os.Exit(1)
		}
		elapsed := time.Since(start)

		if err := os.WriteFile(outPath, decompressed, 0644); err != nil {
			fmt.Printf("[Error] Failed to write extracted file: %v\n", err)
			os.Exit(1)
		}

		shaOut := calculateSHA256(decompressed)
		fmt.Printf("[SUCCESS] Restored to '%s'\n", outPath)
		fmt.Printf("   Restored Size: %d bytes\n", len(decompressed))
		fmt.Printf("   SHA256 Check:  %s\n", shaOut)
		fmt.Printf("   Time Elapsed:  %v\n", elapsed)
		return
	}

	flag.Usage()
}
