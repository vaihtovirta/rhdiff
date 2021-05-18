package differ

import (
	"crypto/sha256"
	"hash/adler32"
	"io"
	"log"
)

// Chunk represents a chunk of data with sha sum and the reference to its index in the original input
type Chunk struct {
	Data      []byte
	Offset    int
	StrongSum [sha256.Size]byte
	WeakSum   uint32
}

// Split splits the reader into chunks of specified length
func Split(reader io.Reader, chunkSizeBytes int) []Chunk {
	chunks := make([]Chunk, 0)
	i := 0

	for {
		buf := make([]byte, chunkSizeBytes)
		n, err := reader.Read(buf)
		buf = buf[:n]

		if n == 0 {
			if err == nil {
				continue
			}

			if err == io.EOF {
				break
			}

			log.Fatal(err)
		}

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		chunk := Chunk{
			Data:      buf,
			Offset:    i,
			WeakSum:   adler32.Checksum(buf),
			StrongSum: sha256.Sum256(buf),
		}

		chunks = append(chunks, chunk)

		i += chunkSizeBytes
	}

	return chunks
}
