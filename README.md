# rhdiff

rhdiff is a simple implementation of input comparison tool based on rolling hash algorithm

# How it works

The library uses the [rsync algorithm](https://rsync.samba.org/tech_report/node2.html) to compare chunks of two files in rolling manner.

Given the delta between two input `src` and `dst` needs to be calculated, `differ` will do the following:

## src preparation
1. Split `src` into non-overlapping chunks of size `S`
1. Calculate weak checksum (adler32) and strong checksum (sha256) for each chunk

## comparison with dst
1. `CalculateDelta` function iterates over each overlapping chunk in `dst` using window of size `S`
1. Create a map from `src` with a weak checksum of each chunk as a key 
1. For each overlapping chunk consecutively check checksums and add moving blocks into the change list
	- Calculate weak checksum of the current window
	- If weak checksum is in the map, calculate strong checksum of the current window
	- If strong and weak checksums match, the chunk in `dst` is the same as in `src`
	- Remove the key-value pair of fully matched from the map
	- If offset of src chunk and offset of dst chunk **match**, the chunk wasn't moved
	- If offset of src chunk and offset of dst chunk **don't match**, the chunk was moved

1. Write bytes between matching chunks into a buffer and dump the changes from the buffer as soon as a next matching block is found

1. Finally, iterate over remaining keys in the `src` map and collect deleted chunks

# Usage

Basic comparison of two strings split into one byte chunks

```go
package main

import (
	"fmt"
	"strings"

	"rhdiff/pkg/differ"
)

const chunkSizeBytes = 3
const includeUnchangedChunks = false

func main() {
	src := bytes.NewReader([]byte("abcxyzfoo"))
	dst := bytes.NewReader([]byte("abc12xyzfo"))

	srcChunks := Split(src, chunkSizeBytes)

	changes := CalculateDelta(srcChunks, dst, chunkSizeBytes)

	for _, change := range changes {
		fmt.Printf(
			"operation: %s, from %d; to %d; text: %s\n",
			change.Operation,
			change.From,
			change.To,
			string(change.Bytes),
		)
	}
}
```

Output

```bash
$ go run examples/basic.go

operation: EQUAL, srcOffset 0; dstOffset 0; data: abc
operation: MOVE, srcOffset 3; dstOffset 5; data: xyz
operation: ADD, srcOffset 0; dstOffset 3; data: 12
operation: ADD, srcOffset 0; dstOffset 8; data: fo
operation: DELETE, srcOffset 6; dstOffset 0; data: foo
```

# Test

```
$ go test ./... -cover
```
