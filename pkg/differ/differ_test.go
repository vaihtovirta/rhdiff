package differ

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCases []struct {
	name      string
	chunkSize int
	src       []byte
	dst       []byte
	changes   []Change
}

func TestDiffer_CalculateChanges(t *testing.T) {
	loremSrc, _ := os.ReadFile("testdata/lorem_src.txt")
	loremDst, _ := os.ReadFile("testdata/lorem_dst.txt")

	cases := testCases{
		{
			name:      "empty source",
			chunkSize: 3,
			src:       []byte(""),
			dst:       []byte("abc"),
			changes: []Change{
				{
					Operation: Add,
					SrcOffset: 0,
					DstOffset: 0,
					Data:      []byte("abc"),
				},
			},
		},
		{
			name:      "empty destination",
			chunkSize: 3,
			src:       []byte("abc"),
			dst:       []byte(""),
			changes: []Change{
				{
					Operation: Delete,
					Data:      []byte("abc"),
				},
			},
		},
		{
			name:      "existing chunks have been changed and new chunks have been added",
			chunkSize: 3,
			src:       []byte("abcxyzfoo"),
			dst:       []byte("abc12xyzfo"),
			changes: []Change{
				{
					Operation: Equal,
					Data:      []byte("abc"),
				},
				{
					Operation: Move,
					SrcOffset: 3,
					DstOffset: 5,
					Data:      []byte("xyz"),
				},
				{
					Operation: Add,
					DstOffset: 3,
					Data:      []byte("12"),
				},
				{
					Operation: Add,
					DstOffset: 8,
					Data:      []byte("fo"),
				},
				{
					Operation: Delete,
					SrcOffset: 6,
					Data:      []byte("foo"),
				},
			},
		},
		{
			name:      "existing chunk has been removed",
			chunkSize: 3,
			src:       []byte("abcxyzfoo"),
			dst:       []byte("abcfoo"),
			changes: []Change{
				{
					Operation: Equal,
					Data:      []byte("abc"),
				},
				{
					Operation: Move,
					SrcOffset: 6,
					DstOffset: 3,
					Data:      []byte("foo"),
				},
				{
					Operation: Delete,
					SrcOffset: 3,
					Data:      []byte("xyz"),
				},
			},
		},
		{
			name:      "text files comparison",
			chunkSize: 64,
			src:       loremSrc,
			dst:       loremDst,
			changes: []Change{
				{
					Operation: Move,
					DstOffset: 49,
					Data:      []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla n"),
				},
				{
					Operation: Add,
					Data:      []byte("Proin finibus ullamcorper ante sit amet egestas. "),
				},
				{
					Operation: Move,
					SrcOffset: 64,
					DstOffset: 113,
					Data:      []byte("isl enim, consectetur quis quam consequat, pharetra tempus enim."),
				},
				{
					Operation: Move,
					SrcOffset: 128,
					DstOffset: 177,
					Data:      []byte(" Fusce iaculis libero vitae ipsum accumsan efficitur. Fusce iacu"),
				},
				{
					Operation: Move,
					SrcOffset: 192,
					DstOffset: 241,
					Data:      []byte("lis est et justo sollicitudin, sed porttitor augue sagittis. Mau"),
				},
				{
					Operation: Move,
					SrcOffset: 320,
					DstOffset: 463,
					Data:      []byte("molestie nisl elit, suscipit egestas ex aliquam ac. Donec dignis"),
				},
				{
					Operation: Add,
					DstOffset: 305,
					Data:      []byte("ris aliquam nisl nibh, sed tempus magna venenatis ac. Nulla ex metus, malesuada eget ultricies vel, fermentum quis nisl. Etiam ac venenatis tellus. Curabitur "),
				},
				{
					Operation: Move,
					SrcOffset: 384,
					DstOffset: 527,
					Data:      []byte("sim, mauris nec malesuada pellentesque, ipsum sem porttitor est,"),
				},
				{
					Operation: Move,
					SrcOffset: 448,
					DstOffset: 591,
					Data:      []byte(" quis laoreet urna orci a leo. Cras tincidunt porttitor sapien, "),
				},
				{
					Operation: Move,
					SrcOffset: 640,
					DstOffset: 714,
					Data:      []byte("."),
				},
				{
					Operation: Add,
					DstOffset: 655,
					Data:      []byte("quis cursus metus pulvinar id. Pellentesque nec mollis eros"),
				},
				{
					Operation: Delete,
					SrcOffset: 256,
					Data:      []byte("ris aliquam nisl nibh, sed tempus magna venenatis ac. Curabitur "),
				},
				{
					Operation: Delete,
					SrcOffset: 512,
					Data:      []byte("quis cursus metus pulvinar id. Pellentesque nec mollis eros. Fus"),
				},
				{
					Operation: Delete,
					SrcOffset: 576,
					Data:      []byte("ce sagittis vehicula ligula, nec ullamcorper sapien sagittis non"),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := bytes.NewReader(c.src)
			dst := bytes.NewReader(c.dst)

			srcChunks := Split(src, c.chunkSize)

			changes := CalculateDelta(srcChunks, dst, c.chunkSize)

			assert.Equal(t, len(c.changes), len(changes))

			for i, expectedChange := range c.changes {
				change := changes[i]

				assert.Equal(t, expectedChange.Operation, change.Operation)
				assert.Equal(t, expectedChange.SrcOffset, change.SrcOffset)
				assert.Equal(t, expectedChange.DstOffset, change.DstOffset)
				assert.Equal(t, expectedChange.Data, change.Data)
			}
		})
	}
}
