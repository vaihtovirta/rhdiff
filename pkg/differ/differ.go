package differ

import (
	"bytes"
	"crypto/sha256"
	"hash/adler32"
	"io"
	"log"
	"sort"
)

// CalculateDelta returns list of changes found between source and destination chunk maps
func CalculateDelta(srcChunks []Chunk, dst *bytes.Reader, chunkSizeBytes int) []Change {
	changes := make([]Change, 0)
	windowStart := 0
	unknownChunkBuffer := new(bytes.Buffer)

	srcWeakCheckSumMap := make(map[uint32]Chunk)
	for _, chunk := range srcChunks {
		srcWeakCheckSumMap[chunk.WeakSum] = chunk
	}

	for {
		windowBuf := make([]byte, chunkSizeBytes)
		n, err := dst.ReadAt(windowBuf, int64(windowStart))
		windowBuf = windowBuf[:n]

		if n == 0 {
			if err == nil {
				continue
			}

			if err == io.EOF {
				writeChangesFromUnknownChunks(
					unknownChunkBuffer,
					&changes,
					windowStart,
				)

				break
			}

			log.Fatal(err)
		}

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		weakSum := adler32.Checksum(windowBuf)
		srcChunk, exists := srcWeakCheckSumMap[weakSum]
		if exists && srcChunk.StrongSum == sha256.Sum256(windowBuf) {
			delete(srcWeakCheckSumMap, weakSum)

			operation := Equal
			if srcChunk.Offset != windowStart {
				operation = Move
			}
			change := Change{
				Operation: operation,
				SrcOffset: srcChunk.Offset,
				DstOffset: windowStart,
				Data:      windowBuf,
			}
			changes = append(changes, change)

			unknownChunkBuffer = writeChangesFromUnknownChunks(
				unknownChunkBuffer,
				&changes,
				windowStart,
			)

			windowStart += chunkSizeBytes
		} else {
			unknownChunkBuffer.WriteByte(windowBuf[0])

			windowStart++
		}
	}

	writeChangesFromMissingChunks(srcWeakCheckSumMap, &changes)

	return changes
}

func writeChangesFromUnknownChunks(buffer *bytes.Buffer, changes *[]Change, offset int) *bytes.Buffer {
	if buffer.Len() > 0 {
		change := Change{
			Operation: Add,
			DstOffset: offset - buffer.Len(),
			Data:      buffer.Bytes(),
		}

		*changes = append(*changes, change)

		return new(bytes.Buffer)
	}

	return buffer
}

func writeChangesFromMissingChunks(srcWeakCheckSumMap map[uint32]Chunk, changes *[]Change) {
	srcWeakSumKeys := make([]uint32, 0)
	for k := range srcWeakCheckSumMap {
		srcWeakSumKeys = append(srcWeakSumKeys, k)
	}

	sort.Slice(
		srcWeakSumKeys,
		func(i, j int) bool { return srcWeakSumKeys[i] > srcWeakSumKeys[j] },
	)

	for _, key := range srcWeakSumKeys {
		chunk := srcWeakCheckSumMap[key]
		change := Change{
			Operation: Delete,
			SrcOffset: chunk.Offset,
			Data:      chunk.Data,
		}

		*changes = append(*changes, change)
	}
}
