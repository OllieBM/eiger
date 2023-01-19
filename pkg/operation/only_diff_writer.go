package operation

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

type ODiffWriter interface {
	DiffWriter
	AddMissingIndexes(indexes []uint64)
}

// will write out a different format
// where any blocks matching their offset will not be
// written out, only changes will be written, so if a block is
// not written at its index, it will be tagged as a remove '-'
type onlyDifferencesWriter struct {
	writer      io.Writer  // not sure if we should include a 'closer'
	mu          sync.Mutex // so we can flush the buffer
	prevOp      *Operation // the last operation, is not written out until there is another operation of it is flushed
	deltaBytes  uint64     // number of additional bytes read
	blockOffset uint64
	//bRead       uint64 // bytes read
}

func NewOnlyDiffWriter(out io.Writer) ODiffWriter {
	return &onlyDifferencesWriter{
		writer:      out,
		prevOp:      nil,
		deltaBytes:  0,
		blockOffset: 0,
	}
}

// old |1234|abcd|5678  [we no longer have abcd]
// new |1234|5678 		[we want to say block0,block2]
// if they are contigious we can create a start, end,
// easy if contigous blocks, we can use two values, begin, and end or length
// if no contigious, we may need to use an []uint32 array  or a []uint64 ->
// each operation could be its own
// 0, 1 increments of chunk count
func (w *onlyDifferencesWriter) AddMatch(index uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	fmt.Println("MATCH AT ", index, w.blockOffset)
	w.flush() // flush on every add right now, since we don't concat them

	// if w.blockOffset == index {
	// 	// this is a match at the correct index. so we ignore it
	// 	return
	// } else
	if w.blockOffset != index {
		// check the difference
		if index < w.blockOffset {
			// either a duplicate block used again or a missing block moved later
			// treat both as an 'ADD'
			w.prevOp = &Operation{operation: OpMatch, blockIndex: index, offset: w.blockOffset}
			return // don't increment the expected block
		} else if index > w.blockOffset {
			// there were some blocks missing!
			// add a remove
			// TODO: w.remove() func
			for i := w.blockOffset; i < index; i++ {
				w.flush()
				fmt.Println("MISSING INDEX ", index, w.blockOffset)
				w.prevOp = &Operation{operation: OpRemoval, blockIndex: i}
			}
			//w.blockOffset = index
		}
		//w.blockOffset = index
		// index == w.blockOffset wont be hit
	} else {
		fmt.Println("MATCHING INDEX ", index, w.blockOffset)
	}

	w.blockOffset = index + 1 // expect the next index
	// its a  match at a different offset

}
func (w *onlyDifferencesWriter) AddMissingIndexes(indexes []uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	sort.Slice(indexes, func(i, j int) bool { return indexes[i] < indexes[j] })

	w.flush()
	for _, i := range indexes {
		if i < w.blockOffset {
			// ignore
			continue
		}

		w.flush()
		w.prevOp = &Operation{operation: OpRemoval, blockIndex: i}
	}

}

// don't mess with indexes on additional characters
func (w *onlyDifferencesWriter) AddMiss(b byte) {
	w.mu.Lock()
	defer w.mu.Unlock()

	fmt.Println("miss at ", w.blockOffset)
	if w.prevOp == nil || w.prevOp.operation != OpMiss {
		// last op was different, flush it
		w.flush()
		// create new op
		w.prevOp = &Operation{operation: OpMiss, data: []byte{b}, offset: w.blockOffset}
		return
	}
	w.prevOp.data = append(w.prevOp.data, b)

}

func (w *onlyDifferencesWriter) Flush() (err error) {

	w.mu.Lock()
	defer w.mu.Unlock()
	w.flush()
	return nil
}

func (w *onlyDifferencesWriter) flush() (err error) {

	if w.prevOp != nil {
		if w.prevOp.operation == OpMiss {
			_, err = fmt.Fprintf(w.writer, "+ @%d %d %s\n", w.prevOp.offset, len(w.prevOp.data), string(w.prevOp.data))
		} else if w.prevOp.operation == OpMatch {
			_, err = fmt.Fprintf(w.writer, "= @%d BLOCK_%d\n", w.prevOp.offset, w.prevOp.blockIndex)
		} else if w.prevOp.operation == OpRemoval {
			_, err = fmt.Fprintf(w.writer, "- BLOCK_%d\n", w.prevOp.blockIndex)
		}
		if err != nil {
			return err
		}
	}
	w.prevOp = nil
	return nil
}

// func (w *onlyDifferencesWriter) Close() error {
// 	if w.closed {
// 		return ErrClosed
// 	}
// 	w.mu.Lock()
// 	defer w.mu.Unlock()
// 	if err := w.flush(); err != nil {
// 		return err
// 	}
// 	return w.writer.Close()
// }
