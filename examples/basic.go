package main

import (
	"bytes"
	"fmt"

	"rhdiff/pkg/differ"
)

const chunkSizeBytes = 3

func main() {
	src := bytes.NewReader([]byte("abcxyzfoo"))
	dst := bytes.NewReader([]byte("abc12xyzfo"))

	srcChunks := differ.Split(src, chunkSizeBytes)

	changes := differ.CalculateDelta(srcChunks, dst, chunkSizeBytes)

	for _, change := range changes {
		fmt.Printf(
			"operation: %s, srcOffset %d; dstOffset %d; data: %s\n",
			change.Operation,
			change.SrcOffset,
			change.DstOffset,
			string(change.Data),
		)
	}
}
