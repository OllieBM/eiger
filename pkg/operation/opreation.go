package operation

const (
	OpMatch OpType = iota // reference
	OpMiss                // addition
	OpRemoval
)

type Operation struct {
	operation  OpType
	data       []byte
	blockIndex uint64
	offset     uint64 // where an add/reference is located in delta target position
}
