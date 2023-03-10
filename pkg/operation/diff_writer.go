package operation

//go:generate mockgen -destination=../mocks/mock_diff_writer.go -package=mocks . DiffWriter

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
)

const (
	insertQuery = `
INSERT INTO users (age, email, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING id
`
)

var ErrClosed = errors.New("DiffWriter is closed")

// DiffWriter is an interface for writing operations to a file
// depending on they are additions or removals
type DiffWriter interface {
	// A block was matched at index of a Signature
	AddMatch(blockIndex uint64)
	// A miss of byte b was detected, this means the byte needs to be added to source
	AddMiss(b byte)
	//  Some blocks were not referenced at all, so an inplace diff writer may need to add tail removal
	AddMissingIndexes(indexes []uint64)
	// flush any currently active operations
	Flush() error
}

// diffWriter will write out a format
// where any blocks matching their offset will not be (and an delta/diff applier would need to )
// written out, only changes will be written, so if a block is
// not written at its index, it will be tagged as a remove operation('-')
type diffWriter struct {
	writer      io.Writer  // not sure if we should include a 'closer'
	mu          sync.Mutex // so we can flush the buffer
	prevOp      Operation  // the last operation, is not written out until there is another operation of it is flushed
	deltaBytes  uint64     // number of additional bytes read
	blockOffset uint64     // the last matched block index, used to keep misses in order

}

// NewDiffWriter will return a writer that implements the diffWriterInterface
func NewDiffWriter(out io.Writer) DiffWriter {
	return &diffWriter{
		writer:      out,
		prevOp:      nil,
		deltaBytes:  0,
		blockOffset: 0,
	}
}

// old |1234|abcd|5678  [we no longer have abcd]
// new |1234|5678 		[we want to provide instructions to ONLY remove CHUNK_1 (abcd)
// which should be close to the minimal amount of data we need to transfer
// Warning!  this is opinionated and will consider a block
// being moved to earlier in a File as
// every expected block from the curent offset to the moved block as
// having been removed
func (w *diffWriter) AddMatch(index uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.prevOp != nil && w.prevOp.Operation() != OpMatch {
		w.flush()
	}

	// if this wasn't the expected block
	if w.blockOffset != index {

		// this is a block that should have appeared earlier
		if index < w.blockOffset {

			// is this a consecutive match
			if w.prevOp != nil && w.prevOp.Operation() == OpMatch {
				p, ok := w.prevOp.(*OperationMatch)
				if ok && p.chunkOffset == index {
					p.chunkOffset = index + 1
					return
				}
			}

			w.flush()
			// either a duplicate block used again or a missing block moved later
			// treat both as an 'ADD'
			w.prevOp = &OperationMatch{chunkOffset: w.blockOffset, chunkStart: index, chunkEnd: index + 1}

			// don't increment the expected block
			return
		} else if index > w.blockOffset {
			w.flush()
			// there were some blocks missing!
			// add a remove for each instance
			// TODO: w.remove() func
			w.addRemoval(&OperationRemoval{chunkStart: w.blockOffset, chunkEnd: index})
		}
	}

	w.blockOffset = index + 1 // expect the next index
}

// AddMissingIndexes add an instance of a removal for each index
// this can be used in cases wh
func (w *diffWriter) AddMissingIndexes(indexes []uint64) {
	if len(indexes) == 0 {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	sort.Slice(indexes, func(i, j int) bool { return indexes[i] < indexes[j] })
	for _, indx := range indexes {
		w.addRemoval(&OperationRemoval{chunkStart: indx, chunkEnd: indx + 1})
	}
}
func (w *diffWriter) addRemoval(r *OperationRemoval) {

	// ignore indexes < blockOffset
	if r.chunkStart < w.blockOffset {
		return
	}

	if w.prevOp != nil && w.prevOp.Operation() == OpRemoval {
		p, ok := w.prevOp.(*OperationRemoval)
		if ok && p.chunkEnd == r.chunkStart {
			p.chunkEnd = r.chunkEnd
			return
		}
	}
	w.flush()
	w.prevOp = r
}

// don't mess with indexes on additional characters
func (w *diffWriter) AddMiss(b byte) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.prevOp != nil && w.prevOp.Operation() == OpMiss {
		p, _ := w.prevOp.(*OperationMiss)
		if p.chunkOffset == w.blockOffset {
			p.data = append(p.data, b)
			return
		}

	}
	// last op was different,  or non consecutive flush it
	w.flush() // send last op
	// create new op
	w.prevOp = &OperationMiss{data: []byte{b}, chunkOffset: w.blockOffset}

}

func (w *diffWriter) Flush() (err error) {

	w.mu.Lock()
	defer w.mu.Unlock()
	w.flush()
	return nil
}

func (w *diffWriter) flush() (err error) {

	if w.prevOp != nil {
		_, err = fmt.Fprintf(w.writer, w.prevOp.String())
		if err != nil {
			return err
		}
	}
	w.prevOp = nil
	return nil
}
