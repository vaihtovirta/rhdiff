package differ

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffer_CalculateChanges(t *testing.T) {
	t.Run("empty source", func(t *testing.T) {
		src := strings.NewReader("")
		dst := strings.NewReader("x")

		srcChunkMap := NewChunksMapFromReader(src, 1)
		dstChunkMap := NewChunksMapFromReader(dst, 1)

		changes := CalculateDelta(srcChunkMap, dstChunkMap)
		sort.Sort(changes)

		assert.Equal(t, 1, len(changes))

		changelist := []Change{
			{
				Operation: Add,
				To:        0,
				Bytes:     []byte("x"),
			},
		}

		for i, expectedChange := range changelist {
			change := changes[i]

			assert.Equal(t, expectedChange.Operation, change.Operation)
			assert.Equal(t, expectedChange.From, change.From)
			assert.Equal(t, expectedChange.To, change.To)
			assert.Equal(t, expectedChange.Bytes, change.Bytes)
		}
	})

	t.Run("empty destination", func(t *testing.T) {
		src := strings.NewReader("x")
		dst := strings.NewReader("")

		srcChunkMap := NewChunksMapFromReader(src, 1)
		dstChunkMap := NewChunksMapFromReader(dst, 1)

		changes := CalculateDelta(srcChunkMap, dstChunkMap)
		sort.Sort(changes)

		assert.Equal(t, 1, len(changes))

		changelist := []Change{
			{
				Operation: Delete,
				From:      0,
				Bytes:     []byte("x"),
			},
		}

		for i, expectedChange := range changelist {
			change := changes[i]

			assert.Equal(t, expectedChange.Operation, change.Operation)
			assert.Equal(t, expectedChange.From, change.From)
			assert.Equal(t, expectedChange.To, change.To)
			assert.Equal(t, expectedChange.Bytes, change.Bytes)
		}
	})

	t.Run("existing chunks have been changed and new chunks have been added", func(t *testing.T) {
		src := strings.NewReader("abc")
		dst := strings.NewReader("xbcd")

		srcChunkMap := NewChunksMapFromReader(src, 1)
		dstChunkMap := NewChunksMapFromReader(dst, 1)

		changes := CalculateDelta(srcChunkMap, dstChunkMap)
		sort.Sort(changes)

		assert.Equal(t, 3, len(changes))

		changelist := []Change{
			{
				Operation: Delete,
				From:      0,
				Bytes:     []byte("a"),
			},
			{
				Operation: Add,
				To:        3,
				Bytes:     []byte("d"),
			},
			{
				Operation: Add,
				To:        0,
				Bytes:     []byte("x"),
			},
		}

		for i, expectedChange := range changelist {
			change := changes[i]

			assert.Equal(t, expectedChange.Operation, change.Operation)
			assert.Equal(t, expectedChange.From, change.From)
			assert.Equal(t, expectedChange.To, change.To)
			assert.Equal(t, expectedChange.Bytes, change.Bytes)
		}
	})

	t.Run("existing chunk has been removed", func(t *testing.T) {
		src := strings.NewReader("abc")
		dst := strings.NewReader("ac")

		srcChunkMap := NewChunksMapFromReader(src, 1)
		dstChunkMap := NewChunksMapFromReader(dst, 1)

		changes := CalculateDelta(srcChunkMap, dstChunkMap)
		sort.Sort(changes)

		assert.Equal(t, 2, len(changes))

		changelist := []Change{
			{
				Operation: Delete,
				From:      1,
				Bytes:     []byte("b"),
			},
			{
				Operation: Move,
				From:      2,
				To:        1,
				Bytes:     []byte("c"),
			},
		}

		for i, expectedChange := range changelist {
			change := changes[i]

			assert.Equal(t, expectedChange.Operation, change.Operation)
			assert.Equal(t, expectedChange.From, change.From)
			assert.Equal(t, expectedChange.To, change.To)
			assert.Equal(t, expectedChange.Bytes, change.Bytes)
		}
	})

	t.Run("new chunks have been inserted between existing ones", func(t *testing.T) {
		src := strings.NewReader("abc")
		dst := strings.NewReader("axybkc")

		srcChunkMap := NewChunksMapFromReader(src, 1)
		dstChunkMap := NewChunksMapFromReader(dst, 1)

		changes := CalculateDelta(srcChunkMap, dstChunkMap)
		sort.Sort(changes)

		assert.Equal(t, 5, len(changes))

		changelist := []Change{
			{
				Operation: Move,
				From:      1,
				To:        3,
				Bytes:     []byte("b"),
			},
			{
				Operation: Move,
				From:      2,
				To:        5,
				Bytes:     []byte("c"),
			},
			{
				Operation: Add,
				To:        4,
				Bytes:     []byte("k"),
			},
			{
				Operation: Add,
				To:        1,
				Bytes:     []byte("x"),
			},
			{
				Operation: Add,
				To:        2,
				Bytes:     []byte("y"),
			},
		}

		for i, expectedChange := range changelist {
			change := changes[i]

			assert.Equal(t, expectedChange.Operation, change.Operation)
			assert.Equal(t, expectedChange.From, change.From)
			assert.Equal(t, expectedChange.To, change.To)
			assert.Equal(t, expectedChange.Bytes, change.Bytes)
		}
	})

	t.Run("file diff with arbitrary chunk size", func(t *testing.T) {
		t.Skip()
		threeKbs := 3 * 1024
		srcFile, _ := os.Open("testdata/lorem_15kb_src.txt")
		dstFile, _ := os.Open("testdata/lorem_15kb_dst.txt")

		srcChunkMap := NewChunksMapFromReader(srcFile, threeKbs)
		dstChunkMap := NewChunksMapFromReader(dstFile, threeKbs)

		changes := CalculateDelta(srcChunkMap, dstChunkMap)
		sort.Sort(changes)

		assert.Equal(t, 14, len(changes))

		assertions := []struct {
			operation OperationType
			text      string
		}{
			{Delete, "Nulla urna eros, sodales a dui quis, fermentum lectus"},
			{Delete, "Ut tincidunt ligula in tellus venenatis matt"},
			{Add, "Maecenas convallis rhoncus sapien posuere ru"},
			{Delete, "Suspendisse lacinia, lectus non semper maximus, massa dolor commodo orci"},
			{Delete, "Cras arcu est, aliquet et elit sit amet"},
			{Add, "Nunc vulputate varius dui"},
			{Delete, "Vivamus rhoncus nulla id augue vene"},
			{Add, "Aenean lobortis, felis sed pellentesque pretium"},
			{Delete, "In sit amet leo et urna dictum pretium"},
			{Delete, "Ut condimentum nibh massa"},
			{Add, "Sed blandit, ipsum ut suscipit imperdiet"},
			{Add, "Mauris vitae ante vitae tortor hendrerit interd"},
			{Add, "In hac habitasse platea dictumst"},
			{Add, "Fusce sed mauris sit amet ipsum hendrerit di"},
		}

		for i, assertion := range assertions {
			change := changes[i]

			assert.Equal(t, assertion.operation, change.Operation)
			assert.Contains(t, string(change.Bytes), assertion.text)
		}
	})
}
