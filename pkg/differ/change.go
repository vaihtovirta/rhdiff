package differ

// OperationType defines enum type for operation types
type OperationType string

const (
	// Add defines enum for ADD operation
	Add OperationType = "ADD"

	// Delete defines enum for DELETE operation
	Delete OperationType = "DELETE"

	// Move defines enum for MOVE operation
	Move OperationType = "MOVE"
)

// Change represents a single change
type Change struct {
	Operation OperationType
	From      int
	To        int
	Bytes     []byte
}

// Delta represents a list of changes
type Delta []Change

// Len returns length of the change slice
func (d Delta) Len() int {
	return len(d)
}

// Less compares two change structs by their string representation
func (d Delta) Less(i, j int) bool {
	return string(d[i].Bytes) < string(d[j].Bytes)
}

// Swap swaps two change structs
func (d Delta) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
