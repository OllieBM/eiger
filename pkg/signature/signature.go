package signature

import (
	"bufio"
	"hash"
	"io"

	"github.com/olliebm/eiger/rolling_checksum"
	"github.com/rs/zerolog/log"
)

type Block struct {
	Index  int
	Strong []byte
	// no weak, its only part of signature
}

type Signature map[uint32][]Block

func Calculate(in io.Reader, blockSize int, strongHasher hash.Hash) (Signature, error) {

	reader := bufio.NewReader(in)
	buf := make([]byte, 0, blockSize)

	rc := rolling_checksum.New()
	signature := make(Signature)
	eof := false
	block_index := 0
	for !eof {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg("could not read from input")
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

		// one weak hash could correspond to several strong hashes
		// so we can keep a list
		signature[weak] = append(signature[weak], Block{Index: block_index, Strong: strong})

		block_index++

	}
	return signature, nil
}
