package signature

import (
	"bufio"
	"bytes"
	"hash"
	"io"

	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/rs/zerolog/log"
)

type Block struct {
	Index  int
	Strong []byte
	// no weak, its only part of signature
}

// type SignatureObj struct {
// 	hashtable    map[uint32][]Block
// 	strongHasher hash.Hash
// }

type hashtable map[uint32][]Block

type Signature struct {
	m            hashtable
	blockSize    int
	strongHasher hash.Hash
	rc           rolling_checksum.RollingChecksum
}

func New(in io.Reader, blockSize int, strongHasher hash.Hash, rc rolling_checksum.RollingChecksum) (*Signature, error) {

	reader := bufio.NewReader(in)
	buf := make([]byte, blockSize)

	m := make(hashtable)
	eof := false
	block_index := 0
	for !eof {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {

				log.Error().Err(err).Msg("could not read from input")
				return nil, err
			}
			eof = true
		}

		if n == 0 {
			continue
		}

		buf = buf[:n] // trim of waste

		// calculate strong hash
		strongHasher.Reset()
		n, err = strongHasher.Write(buf)
		if err != nil {
			log.Error().Err(err).Msg("could not create strong hash")
		}
		strong := strongHasher.Sum(nil)

		// calculate weak hash
		weak := rc.Calculate(buf)
		log.Debug().Msgf("adding signature: '%s' weak:[%d])", string(buf), weak)
		// one weak hash could correspond to several strong hashes
		// so we can keep a list
		m[weak] = append(m[weak], Block{Index: block_index, Strong: strong})

		block_index++

	}

	return &Signature{
		m:            m,
		blockSize:    blockSize,
		strongHasher: strongHasher,
	}, nil

}

func (s *Signature) FindMatch(weak uint32, buf []byte) (bool, int) {
	s.strongHasher.Reset()
	if hashes, ok := s.m[weak]; ok {
		_, _ = s.strongHasher.Write(buf)
		strong := s.strongHasher.Sum(nil)
		for _, h := range hashes {
			if bytes.Compare(strong, h.Strong) == 0 {
				return true, h.Index
			}
		}
	}
	return false, 0
}

func (s Signature) BlockSize() int {
	return s.blockSize
}

func (s Signature) Hashtable() hashtable {
	return s.m
}
