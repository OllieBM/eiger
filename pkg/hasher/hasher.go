package hasher

import (
	"bufio"
	"hash"
	"io"

	"github.com/OllieBM/eiger/pkg/rolling_checksum"
)

// TODO:
// consider this to make passing a hash.Hash around simpler
// and removing the .sum(notnil)
type Hasher interface {
	Strong([]byte) []byte
	Weak([]byte) uint32 // rolling
	//Reset(io.Reader, uint32)
}

type hasher struct {
	hasher hash.Hash
	//rolling rolling_checksum.RollingChecksum
}

// convenience type
type RollingHasher interface {
	Roll() (uint32, error)      // roll the checksum, consumes 1 byte or triggers a recalculation
	Hash() uint32               // get current hash value
	Calculate() (uint32, error) // recalculate the current hash, reads byte length
}

type rollingHasher struct {
	reader    *bufio.Reader
	hash      uint32
	chunkSize uint32
	buffer    []byte // TODO: use ring buffer
	//dirty     bool   // need to recalculate hash
	rc rolling_checksum.RollingChecksum
}

// New Creates a rolling hasher
func New(in io.Reader, chunkSize uint32) RollingHasher {
	return &rollingHasher{
		reader:    bufio.NewReader(in),
		hash:      0,
		chunkSize: chunkSize,
		buffer:    make([]byte, chunkSize),
		rc:        rolling_checksum.New(),
	}
}

// Roll will try to read 1 more byte from the underlying io.Reader
// if we were able to successfully read a byte we use our
// rolling checksum to cycle that value
func (r *rollingHasher) Roll() (uint32, error) {
	// read a single bit
	prv := r.buffer[0]
	r.buffer = r.buffer[1:]
	b, err := r.reader.ReadByte()
	if err != nil {
		if err != io.EOF {
			return 0, err
		}
		// we didn't read a byte
		// recalculate buffer
		r.hash = r.rc.Calculate(r.buffer)
		if len(r.buffer) > 0 {
			// we ignore eof until the entire buffer is empty
			err = nil
		}
	} else {
		r.hash = r.rc.Roll(prv, b)
		r.buffer = append(r.buffer, b)
	}
	return r.hash, err

}

// Calculate will read N bytes where N <= Chunk Size
// and create a hash based on the returned N bytes
// if the underlying reader is empty  io.EOF is returned
func (r *rollingHasher) Calculate() (uint32, error) {

	// read chunk size bytes and then calculate a new hash
	n, err := r.reader.Read(r.buffer)
	if err != nil {
		if err != io.EOF {
			return 0, err
		}
		// we got a io.EOF
		// if we read enough
		if n == 0 {
			// read all elements from buffer
			return 0, io.EOF
		}
	}
	// trim excess if buffer read less than chunk size
	r.buffer = r.buffer[:n]
	r.hash = r.rc.Calculate(r.buffer)
	return r.hash, nil

}
func (r *rollingHasher) Hash() uint32 {
	return r.Hash()
}
