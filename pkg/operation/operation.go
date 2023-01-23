package operation

import "fmt"

type OpType uint8

const (
	OpMatch OpType = iota // match
	OpMiss                // addition
	OpRemoval
)

// Operation interface providing a value for checking which subtype it is
// and a String method for printing out operation
type Operation interface {
	Operation() OpType
	String() string
}

//////////////////
// OpMatch		//
//////////////////
type OperationMatch struct {
	chunkOffset uint64 // where an add/match is located in the source file
	chunkStart  uint64 // which chunk is matched
	chunkEnd    uint64 // which was the last chunk matched consecutively
}

func (o OperationMatch) Operation() OpType {
	return OpMatch
}
func (o OperationMatch) String() string {
	return fmt.Sprintf("= @%d CHUNK_%d\n", o.chunkOffset, o.chunkStart)
}

//////////////////
// OpMiss		//
//////////////////
type OperationMiss struct {
	chunkOffset uint64 // where an add/match is located in delta target position
	data        []byte // characters to add to the target before the chunks at chunkOffset
}

func (o OperationMiss) Operation() OpType {
	return OpMiss
}
func (o OperationMiss) String() string {
	return fmt.Sprintf("+ @%d %d %s\n", o.chunkOffset, len(o.data), string(o.data))
}

//////////////////
// OpRemove		//
//////////////////
type OperationRemoval struct {
	chunkStart uint64 // which chunk a remove begins
	chunkEnd   uint64 // which chunk a remove ends
}

func (o OperationRemoval) Operation() OpType {
	return OpRemoval
}

func (o OperationRemoval) String() string {
	return fmt.Sprintf("- CHUNK_%d...CHUNK_%d\n", o.chunkStart, o.chunkEnd)
}
