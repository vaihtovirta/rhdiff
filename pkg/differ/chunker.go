package differ

import (
	"bufio"
	"crypto/sha256"
	"io"
	"log"
)

// Chunk represents a chunk of data with sha sum and the reference to its index in the original input
type Chunk struct {
	Bytes []byte
	Index int
	Sum   [sha256.Size]byte
}

// ChunkMap represents map between chunk's sum and the Chunk itself
type ChunkMap map[[sha256.Size]byte]Chunk

// NewChunksMap creates a new ChunkMap from the given byte chunks
func NewChunksMap(byteChunks [][]byte) ChunkMap {
	chunksMap := make(ChunkMap)

	for i, byteChunk := range byteChunks {
		sum := sha256.Sum256(byteChunk)

		chunksMap[sum] = Chunk{
			Bytes: byteChunk,
			Index: i,
			Sum:   sum,
		}
	}

	return chunksMap
}

// NewChunksMap creates a new ChunkMap from the given IO reader and chunk size in bytes
func NewChunksMapFromReader(r io.Reader, chunkSizeBytes int) ChunkMap {
	reader := bufio.NewReader(r)
	chunks := Split(reader, chunkSizeBytes)

	return NewChunksMap(chunks)
}

// Split splits the input into byte chunks of given size
func Split(reader io.Reader, chunkSizeBytes int) [][]byte {
	chunks := make([][]byte, 0)

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

		chunks = append(chunks, buf)
	}

	return chunks
}
