// Package ipfs implements the chunking, encryption, and CID logic of Section 3.2 / Algorithm 2.
package ipfs

import "bytes"

const DefaultChunkSize = 256 * 1024 // 256 KiB

// Chunk splits data into fixed-size chunks (Algorithm 5, line 4).
func Chunk(data []byte, size int) [][]byte {
	if size <= 0 { size = DefaultChunkSize }
	var out [][]byte
	r := bytes.NewReader(data)
	buf := make([]byte, size)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			cp := make([]byte, n)
			copy(cp, buf[:n])
			out = append(out, cp)
		}
		if err != nil { break }
	}
	return out
}