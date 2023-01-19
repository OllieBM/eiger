package delta

import (
	"bufio"
	"errors"
	"io"

	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidSignature = errors.New("Nil signature passed to delta.Calculate")
	ErrInvalidOpWriter  = errors.New("Nil OpWriter passed to delta.Calculate")
)

func Calculate(in io.Reader, sig *signature.Signature, out *operation.OpWriter) error {
	if sig == nil {
		log.Error().Err(ErrInvalidSignature)
		return ErrInvalidSignature
	}
	if out == nil {
		log.Error().Err(ErrInvalidOpWriter)
		return ErrInvalidOpWriter
	}

	reader := bufio.NewReader(in)
	r := rolling_checksum.New()

	eof := false
	rolling := false
	buf := make([]byte, sig.BlockSize())
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
		log.Debug().Msgf("searching for match: '%s' weak:[%d])", string(buf), weak)
		match, indx := sig.FindMatch(weak, buf)
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
					log.Debug().Msgf("searching for match: '%s' weak:[%d])", string(buf), weak)
					match, indx := sig.FindMatch(weak, buf)
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
func Calculate2(in io.Reader, sig *signature.Signature, out operation.DiffWriter) error {
	if sig == nil {
		log.Error().Err(ErrInvalidSignature)
		return ErrInvalidSignature
	}
	if out == nil {
		log.Error().Err(ErrInvalidOpWriter)
		return ErrInvalidOpWriter
	}

	reader := bufio.NewReader(in)
	r := rolling_checksum.New()

	chunk := make([]byte, sig.BlockSize())
	//one := make([]byte, 1)
	//buf := chunk // shadow the underlying buffer

	rolling := false

	var n int
	var b, prevC byte
	var weak uint32

	var err error

	// read either [1]byte or [blockSize]byte from reader
	// until an error occurs

	// read buf amount of bytes,
	// it is either a buffer of 1 or blocksize
	rolling = false
	for err == nil {
		if !rolling {
			n, err = reader.Read(chunk)
			weak = r.Calculate(chunk)
			chunk = chunk[:n]
			log.Debug().Msgf("read %s", chunk)
			if n == 0 {
				break
			}

		} else {

			// b, err = reader.ReadByte()
			// if err != nil {
			// 	break
			// }
			// chunk = append(chunk[1:], b)
			// weak = r.Roll(prevC, b)

			b, err = reader.ReadByte()
			if err != nil {
				break
			}
			chunk = append(chunk[1:], b)
			weak = r.Roll(prevC, b)
		}

		match, indx := sig.FindMatch(weak, chunk)
		if match {
			log.Debug().Msgf("Match for '%s' weak[%d]", chunk, weak)
			rolling = false
			out.AddMatch(uint64(indx))
		} else {
			rolling = true
			prevC = chunk[0]
			log.Debug().Msgf("Miss for '%s' weak[%d] adding %s", string(chunk), weak, string(prevC))
			out.AddMiss(prevC)
		}
	}

	// missing chunks :
	if err != nil {
		if err != io.EOF {
			log.Error().Err(err)
			return err
		}
	}
	return nil
}

// calculate2 tracking any unused chunks
func Calculate3(in io.Reader, sig *signature.Signature, out operation.ODiffWriter) error {
	if sig == nil {
		log.Error().Err(ErrInvalidSignature)
		return ErrInvalidSignature
	}
	if out == nil {
		log.Error().Err(ErrInvalidOpWriter)
		return ErrInvalidOpWriter
	}

	// matched chunks :=
	matched := make(map[uint64]struct{})

	reader := bufio.NewReader(in)
	r := rolling_checksum.New()

	chunk := make([]byte, sig.BlockSize())
	//one := make([]byte, 1)
	//buf := chunk // shadow the underlying buffer

	rolling := false

	var n int
	var b, prevC byte
	var weak uint32

	var err error

	// read either [1]byte or [blockSize]byte from reader
	// until an error occurs

	// read buf amount of bytes,
	// it is either a buffer of 1 or blocksize
	rolling = false
	for err == nil {
		if !rolling {
			n, err = reader.Read(chunk)
			chunk = chunk[:n]
			weak = r.Calculate(chunk)
			log.Debug().Msgf("read %s", chunk)
			if n == 0 {
				break
			}

		} else {

			b, err = reader.ReadByte()
			if err != nil {
				break
			}
			chunk = append(chunk[1:], b)
			weak = r.Roll(prevC, b)
		}

		match, indx := sig.FindMatch(weak, chunk)
		if match {
			log.Debug().Msgf("Match for '%s' weak[%d]", chunk, weak)
			rolling = false
			out.AddMatch(uint64(indx))
			matched[uint64(indx)] = struct{}{}

		} else {
			rolling = true
			prevC = chunk[0]
			log.Debug().Msgf("Miss for '%s' weak[%d] adding %s", string(chunk), weak, string(prevC))
			out.AddMiss(prevC)
		}
	}

	// write out missing characters if there are some left
	if rolling {
		chunk = chunk[1:]
		for len(chunk) > 0 {
			out.AddMiss(chunk[0])
			chunk = chunk[1:]
		}
	}
	// write out any unmatched chunks
	missing := make([]uint64, 0)
	for _, v := range sig.Hashtable() {
		for _, block := range v {

			if _, ok := matched[uint64(block.Index)]; !ok {
				missing = append(missing, uint64(block.Index))
			}
		}
	}
	out.AddMissingIndexes(missing)

	// missing chunks :
	if err != nil {
		if err != io.EOF {
			log.Error().Err(err)
			return err
		}
	}

	return nil
}
