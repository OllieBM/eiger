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

func Calculate(in io.Reader, sig signature.Signature, hasher hash.Hash, blockSize uint64, out operation.OpWriter) error {

	reader := bufio.NewReader(in)

	opW := operation.OpWriter{}
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
			weak = r.Calculate(buf[:n])
		} else {
			// read one byte
			lastByte := buf[0]
			var b byte
			b, err = reader.ReadByte()
			buf = append(buf[1:blockSize], b)
			weak = r.Roll(lastByte, b)
		}
		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg("failed to read delta input")
				return err
			}
			eof = true
		}

		// look for a match in signature
		match, indx := FindMatch(weak, buf, hasher, sig)
		if match {
			opW.AddMatch(uint64(indx))
			rolling = false
			continue
		}
		if !eof {
			// miss
			opW.AddMiss(buf[0])
		} else {
			// add all remaining
			for _, c := range buf {
				opW.AddMiss(c)
			}
		}
	}

	return nil
}

func FindMatch(weak uint32, buf []byte, hasher hash.Hash, sig signature.Signature) (bool, int) {
	if hashes, ok := sig[weak]; ok {
		hasher.Reset()
		strong := hasher.Sum(buf)
		for _, h := range hashes {
			if bytes.Compare(strong, h.Strong) == 0 {
				return true, h.Index
			}
		}
	}
	return false, 0
}
