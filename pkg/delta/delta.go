package delta

import (
	"bufio"
	"bytes"
	"hash"
	"io"

	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/rs/zerolog/log"
)

func Calculate(in io.Reader, sig signature.Signature, hasher hash.Hash, blockSize uint64, out *operation.OpWriter) error {

	reader := bufio.NewReader(in)
	r := rolling_checksum.New()

	eof := false
	rolling := false
	buf := make([]byte, blockSize)
	var n int
	var err error
	for !eof {

		var weak uint32
		if !rolling {
			// read buf
			n, err = reader.Read(buf)
			// if n == blockSize{ // sHould we skip/}
			buf = buf[:n]

			weak = r.Calculate(buf)
		} else {
			// read one byte
			lastByte := buf[0]
			buf = buf[1:] // remove the element from buffer
			var b byte
			b, err = reader.ReadByte()
			if err == nil {
				// read will return a default byte
				// and an error if something cannot be read
				buf = append(buf, b)
				weak = r.Roll(lastByte, b)
			} else {
				// we have changes the block size
				weak = r.Calculate(buf)
			}

		}
		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg("failed to read delta input")
				return err
			}
			eof = true
		}
		if len(buf) == 0 {
			continue
		}

		// look for a match in signature
		log.Debug().Msgf("searchig for match: '%s' weak:[%d])", string(buf), weak)
		match, indx := FindMatch(weak, buf, hasher, sig)
		if match {
			log.Debug().Msgf("Match for '%s' weak[%d]", buf, weak)
			out.AddMatch(uint64(indx))
			rolling = false
			continue
		} else {
			log.Debug().Msgf("Miss for '%s' weak[%d]", buf, weak)
			out.AddMiss(buf[0])
			rolling = true
			if eof {
				// flush the last characters in the buffer
				// and try matches or add misses
				// TODO:
				// we could just add everything as a miss
				buf = buf[1:]
				for len(buf) > 0 {
					weak := r.Calculate(buf)
					log.Debug().Msgf("searchig for match: '%s' weak:[%d])", string(buf), weak)
					match, indx := FindMatch(weak, buf, hasher, sig)
					if match {
						log.Debug().Msgf("Match for '%s' %d", buf, weak)
						out.AddMatch(uint64(indx))
						buf = nil
					} else {

						log.Debug().Msgf("Miss for '%s' %d", buf, weak)
						out.AddMiss(buf[0])
						buf = buf[1:]
					}
				}
			}

			// if !eof {
			// 	// we now roll
			// 	rolling = true
			// 	// miss
			// 	out.AddMiss(buf[0])

			// } else {
			// 	// now we roll  and try to find a match

			// 	// add all remaining
			// 	for _, c := range buf {
			// 		out.AddMiss(c)
			// 	}
			// }
		}
	}

	return nil
}

// Refactor
// Calculate2 will generate a delta between 'in' and the signature file
// in should be the 'leader' file and the signature file should be based
// on the 'follower' file.
func Calculate2(in io.Reader, sig signature.Signature, hasher hash.Hash, blockSize uint64, out *operation.OpWriter) error {

	reader := bufio.NewReader(in)
	r := rolling_checksum.New()

	chunk := make([]byte, blockSize)
	one := make([]byte, 1)
	buf := chunk // shadow the underlying buffer

	rolling := false

	var prevC byte
	var weak uint32
	var n int
	var err error
	// read either [1]byte or [blockSize]byte from reader
	// until an error occurs
	for n, err = reader.Read(buf); err == nil; n, err = reader.Read(buf) {
		if !rolling {
			weak = r.Calculate(buf)
		} else {
			weak = r.Roll(prevC, buf[len(buf)-1])
		}

		log.Debug().Msgf("finding match for %s [%d]", buf, weak)
		match, indx := FindMatch(weak, buf, hasher, sig)
		if match {
			rolling = false
			buf = chunk // read chunk
			out.AddMatch(uint64(indx))
		} else {
			rolling = true
			prevC = buf[0]
			buf = one
			out.AddMiss(buf[0])
		}
	}
	// tail
	// TODO: wrap up into a closure or function to streamline code
	if n != 0 {
		for len(buf) > 0 {
			weak = r.Calculate(buf)
			match, indx := FindMatch(weak, buf, hasher, sig)
			if match {
				out.AddMatch(uint64(indx))
				break
			} else {
				prevC = buf[0]
				out.AddMiss(prevC)
				buf = buf[1:]
			}
		}
	}

	if err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}

func FindMatch(weak uint32, buf []byte, hasher hash.Hash, sig signature.Signature) (bool, int) {
	hasher.Reset()
	if hashes, ok := sig[weak]; ok {
		_, _ = hasher.Write(buf)
		strong := hasher.Sum(nil)
		for _, h := range hashes {
			if bytes.Compare(strong, h.Strong) == 0 {
				return true, h.Index
			}
		}
	}
	return false, 0
}
