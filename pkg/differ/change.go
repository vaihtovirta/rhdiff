package differ

// OperationType defines enum type for operation types
type OperationType string

const (
	// Equal defines enum for EQUAL operation
	Equal OperationType = "EQUAL"

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
	SrcOffset int
	DstOffset int
	Data      []byte
}
