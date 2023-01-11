package operation

import (
	"fmt"
	"io"
)

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

// Try and write a diff
// empties the internal buffer
func (w *OpWriter) Output(out io.Writer) error {
	var err error

	// for _, op := range w.ops {
	for len(w.ops) > 0 {
		op := w.ops[0]
		if op.operation == OpMiss {
			_, err = fmt.Fprintf(out, "+ %d %s\n", len(op.data), string(op.data))
		} else if op.operation == OpMatch {
			_, err = fmt.Fprintf(out, "= BLOCK_%d\n", op.blockIndex)
		}
		if err != nil {
			return err
		}
		w.ops = w.ops[1:]
	}
	w.ops = nil
	return nil
}

// Try and write a diff
// empties the internal buffer
func (w *OpWriter) Flush(out io.Writer) error {
	var err error

	// for _, op := range w.ops {
	for len(w.ops) > 0 {
		op := w.ops[0]
		if op.operation == OpMiss {
			_, err = fmt.Fprintf(out, "+ %d %s\n", len(op.data), string(op.data))
		} else if op.operation == OpMatch {
			_, err = fmt.Fprintf(out, "= BLOCK_%d\n", op.blockIndex)
		}
		if err != nil {
			return err
		}
		w.ops = w.ops[1:]
	}
	w.ops = nil
	return nil
}
