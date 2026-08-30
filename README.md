<div align="center">
  <h1>📦 Huffman & LZW Data Compression</h1>
  
  <p><b>A high-performance, zero-dependency Computer Science suite implementing Huffman Coding, LZW (Lempel-Ziv-Welch), Run-Length Encoding (RLE), and Shannon Information Entropy Analysis in pure Go.</b></p>

  <p>
    <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
    <img src="https://img.shields.io/badge/Dependencies-Zero-00C853?style=for-the-badge" />
    <img src="https://img.shields.io/badge/Algorithms-Huffman%20%7C%20LZW%20%7C%20RLE-6B46C1?style=for-the-badge" />
    <img src="https://img.shields.io/badge/Theory-Shannon%20Entropy-FF5F87?style=for-the-badge" />
  </p>
</div>

---

## 🔬 Overview & Scientific Background

Data compression is the science of encoding information using fewer bits than the original representation. This repository provides modular, from-scratch implementations of fundamental loss-less compression algorithms and information-theoretic analyzers with **zero third-party dependencies**.

```text
               +-----------------------------+
               |      Raw Input Stream       |
               +-----------------------------+
                              |
            +-----------------+-----------------+
            |                 |                 |
     [Huffman Tree]       [LZW Dict]          [RLE]
   Variable-length      Dictionary Prefix   Repetitive Run
     Prefix Codes         Substitutions        Sequences
            |                 |                 |
            +-----------------+-----------------+
                              |
               +-----------------------------+
               | Custom BitStream Container  |
               | (Header + Payload + SHA256) |
               +-----------------------------+
```

---

## ⚡ Implemented Algorithms

### 1. Huffman Coding (`pkg/huffman`)
- **Min-Heap Construction:** Uses a binary priority queue to merge the lowest-frequency character nodes in $O(N \log K)$ time.
- **Prefix-Free Codes:** Traverses the binary tree to generate variable-length bit codes (more frequent characters receive shorter bit sequences).
- **Tree Topology Serialization:** Encodes the binary tree structure directly into the archive bitstream using pre-order bit traversal ($0$ for internal node, $1$ + $8\text{-bit symbol}$ for leaf).

### 2. LZW Dictionary Compression (`pkg/lzw`)
- Dynamically constructs a string dictionary on the fly without transmitting the dictionary in the compressed file.
- Encodes recognized string patterns into fixed $12\text{-bit}$ codewords (supporting up to 4096 dynamic phrases).

### 3. Run-Length Encoding (`pkg/rle`)
- Fast linear byte scanner that replaces contiguous identical byte runs with `(byte, count)` pairs. Highly effective for sparse or uncompressed bitmap assets.

### 4. Shannon Entropy Analyzer (`pkg/analyzer`)
- Calculates the fundamental theoretical limit of lossless data compression based on Claude Shannon's Information Theory:
  $$H(X) = -\sum_{i=1}^{n} P(x_i) \log_2 P(x_i)$$
- Determines the exact average minimum bits required per byte symbol and computes the maximum theoretical space savings percentage.

---

## 🛠️ Installation & Building

No external tools or packages are needed. Built entirely with the standard Go toolchain:

```bash
# Clone the repository
git clone https://github.com/amiraliT-88/Huffman-LZW-Compression.git
cd Huffman-LZW-Compression

# Build standalone executable
go build -o datacompression.exe .
```

---

## 💻 CLI Usage Guide

### 1. Analyze Shannon Information Entropy
Calculates the theoretical compression ceiling and symbol frequency distribution:
```bash
./datacompression.exe -analyze -in=sample.txt
```

### 2. Cross-Algorithm Benchmark
Compares the throughput (MB/s) and space savings of all three algorithms side-by-side:
```bash
./datacompression.exe -benchmark -in=sample.txt
```

### 3. Compress a File
Compress using your algorithm of choice (`huffman`, `lzw`, or `rle`):
```bash
# Compress with Huffman Coding (default)
./datacompression.exe -c -algo=huffman -in=document.txt -out=document.txt.dca

# Compress with LZW Dictionary
./datacompression.exe -c -algo=lzw -in=data.bin -out=data.bin.dca
```

### 4. Decompress an Archive
The engine automatically detects the compression algorithm from the binary container header:
```bash
./datacompression.exe -d -in=document.txt.dca -out=restored_document.txt
```

---

## 🧪 Unit Testing & Verification

Run the comprehensive unit test suite to verify lossless round-trip compression across plain text, repetitive streams, all 256-symbol alphabets, and random binary blocks:

```bash
go test -v ./...
```

---

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
