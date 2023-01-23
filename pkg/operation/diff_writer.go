package operation

import "errors"

var ErrClosed = errors.New("DiffWriter is closed")

// DiffWriter is an interface for writing operations to a file
// depending on they are additions or removals
type DiffWriter interface {
	// A block was matched at index of a Signature
	AddMatch(blockIndex uint64)
	// A miss of byte b was detected, this means the byte needs to be added to source
	AddMiss(b byte)
	Flush() error
}
