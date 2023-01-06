package rolling_hash

import (
	"bufio"
	"fmt"
	"os"
)

// A Prime random large prime number
const prime = 1190494759 // selected from http://www.primos.mat.br/2T_en.html
const base = 101

type RollingHasher struct {
	windowSize uint8
	bufferSize uint8
	oldFile    *bufio.Reader
	newFile    *bufio.Reader
	closed     chan struct{}
}

const NullHasher := RollingHasher{}

func NewRollingHashser(oldFile string, newFile string, bufferSize uint8) (RollingHasher, error) {

	// TODO:
	// 1. check files exist
	// open scanners

	result := RollingHasher{}

	file1, err := os.Open(oldFile)
	if err != nil {
		return result, fmt.Errorf("failed to open %s: %w", oldFile, err)
	}
	file2, err := os.Open(newFile)
	if err != nil {
		return result, fmt.Errorf("failed to open %s: %w", newFile, err)
	}

	return RollingHasher{
		windowSize: 20,
		bufferOne:  make([]byte, 0, bufferSize),
		bufferTwo:  make([]byte, 0, bufferSize),
		oldFile:    bufio.NewReader(file1),
		newFile:    bufio.NewReader(file2),
		closed:     make(chan struct{}),
	}, nil
}

func (rh *RollingHasher) ProcessChunks() {
	// read file old file as chunks
	rh.oldFile.Read()
	
}


// 1. read a file create chunk checksums
// 2. rolling checksums for chunks that fail
// rolling checksums are used to compute inter block mis matches using a 
// cache of the original


func 