package operation

type OpType uint8

const (
	OpMatch OpType = iota // reference
	OpMiss                // addition
)

type Operation struct {
	operation  OpType
	data       []byte
	blockIndex uint64
}

type OpWriter struct {
	ops []Operation
}

func (w *OpWriter) AddMatch(index uint64) {
	w.ops = append(w.ops, Operation{operation: OpMatch, blockIndex: index})
}

func (w *OpWriter) AddMiss(b byte) {
	w.addMiss([]byte{b})
}

func (w *OpWriter) addMiss(bytes []byte) {
	n := len(w.ops)
	if n > 0 && w.ops[n-1].operation == OpMiss {
		// extend the miss operation
		w.ops[n-1].data = append(w.ops[n-1].data, bytes...)
		return
	}
	w.ops = append(w.ops, Operation{operation: OpMiss, data: bytes})
}

func (w *OpWriter) Operations() []Operation {
	return w.ops
}
