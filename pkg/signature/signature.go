package signature

import (
	"bufio"
	"strings"

	"github.com/rs/zerolog/log"
)

type signature = uint32

func process(s string, chunkSize uint) {

	buffer := make([]byte, 0, chunkSize)
	reader := bufio.NewReader(strings.NewReader(s))

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			log.Error().Error(err).Msg("error reading buffer")
			break
		}

	}

	for {

		n, err := in.Read(buf)
		if n != 0 {
			// remember that the last block could be < chunkLen
			// and is acceptable
			buf = buf[:n] // trim off garbage

			// strong hash
			// create strong hash
			strongHasher.Reset()
			strongHasher.Write(buf)
			strong := strongHasher.Sum(nil)

			// create rolling hash
			weak, _, _ := rolling_hash.Hash(buf)
			// weak := rolling_hash.rolling_hash(buf, prime, base)

			c := Chunk{
				Index:  index,
				Strong: strong,
				Weak:   weak,
			}
			chunks = append(chunks, c)
		}
		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg("reading file")
				return nil, err // break
			}
			break
		}
	}
	return chunks, nil
}

}
