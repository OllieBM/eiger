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

// Calculate a delta from in based on differences between it and another file defined by signature
// the input source should be the Leader or source of truth, while the signature should be calculated from
// a follower or stale file.
func Calculate(in io.Reader, sig *signature.Signature, rc rolling_checksum.RollingChecksum, out operation.DiffWriter) error {

	if sig == nil {
		log.Error().Err(ErrInvalidSignature).Msg("")
		return ErrInvalidSignature
	}
	if out == nil {
		log.Error().Err(ErrInvalidOpWriter).Msg("")
		return ErrInvalidOpWriter
	}

	reader := bufio.NewReader(in)
	chunk := make([]byte, sig.BlockSize())

	rolling := false

	var n int
	var b, prevC byte
	var weak uint32
	var err error
	// this is to keep track of any tail removals (things left over by the signature a)
	// only needed in diff writers where duplicate sections are not written out
	// (i.e. we only make entries for moves, adds, and deletes, not where things are identical)
	matched := make(map[uint64]struct{})

	// read either 1 byte, of chunk_size amount of bytes, until an error occurs
	rolling = false
	for err == nil {
		if !rolling {
			n, err = reader.Read(chunk)
			chunk = chunk[:n]
			weak = rc.Calculate(chunk)
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
			weak = rc.Roll(prevC, b)
		}

		match, indx := sig.FindMatch(weak, chunk)
		if match {
			log.Debug().Msgf("Match for '%s' weak[%d] with index [%d]", chunk, weak, indx)
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

	// early exit if there was a fatal errors
	if err != nil {
		if err != io.EOF {
			log.Error().Err(err).Msg("")
			return err
		}
	}

	// write out any unmatched chunks from the source signature?
	missing := make([]uint64, 0)
	for _, v := range sig.Hashtable() {
		for _, block := range v {

			if _, ok := matched[uint64(block.Index)]; !ok {
				missing = append(missing, uint64(block.Index))
			}
		}
	}

	// tell the diff writer that these indexes are missing
	out.AddMissingIndexes(missing)

	return nil
}
