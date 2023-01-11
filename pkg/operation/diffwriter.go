package operation

import (
	"fmt"
	"io"
	"sync"
)

// const MaxOpLength =  1 << 16 // maximum amount of data one operation can be
// TODO: cap size of previous operation
// TODO: add in a way for MatchOperations to contain multi-block references (start, end) vs []index

type DiffWriter interface {
	AddMatch(index uint64) // index, length
	AddMiss(b byte)
	//io.WriteCloser
	io.Closer
	//Flush()
}

type customDiffWriter struct {
	writer io.WriteCloser // not sure if we should close it
	mu     sync.Mutex     // so we can flush the buffer
	prevOp *Operation     // the last operation, is not written out until there is another operation of it is flushed
}

func NewDiffWriter(out io.Writer) DiffWriter {
	return &customDiffWriter{}
}

// old |1234|abcd|5678  [we no longer have abcd]
// new |1234|5678 		[we want to say block0,block2]
// if they are contigious we can create a start, end,
// easy if contigous blocks, we can use two values, begin, and end or length
// if no contigious, we may need to use an []uint32 array  or a []uint64 ->
// each operation could be its own
// 0, 1 increments of chunk count
func (w *customDiffWriter) AddMatch(index uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	// if w.prevOp.operation == OpMiss {
	// 	Flush()
	// }
	w.Flush() // flush on every add right now, since we don't concat them
	w.prevOp = &Operation{operation: OpMatch, blockIndex: index}
}

func (w *customDiffWriter) AddMiss(b byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.prevOp.operation != OpMiss {
		// last op was different, flush it
		w.Flush()
		// create new op
		w.prevOp = &Operation{operation: OpMiss, data: []byte{b}}
		return
	}
	w.prevOp.data = append(w.prevOp.data, b)

}

func (w *customDiffWriter) Flush() (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.prevOp != nil {
		if w.prevOp.operation == OpMiss {
			_, err = fmt.Fprintf(w.writer, "+ %d %s\n", len(w.prevOp.data), string(w.prevOp.data))
		} else if w.prevOp.operation == OpMatch {
			_, err = fmt.Fprintf(w.writer, "= BLOCK_%d\n", w.prevOp.blockIndex)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *customDiffWriter) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.writer.Close()
}
