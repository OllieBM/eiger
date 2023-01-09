package delta

import (
	"bufio"
	"io"

	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/rs/zerolog/log"
)

func Calculate(in io.Reader, sig signature.Signature) error {

	reader := bufio.NewReader(in)

	if len(sig) == 0 {
		// we should treat the entire input as a new file
		buf := make([]byte, 0, 4096)
		eof := false
		for !eof {
			, err := reader.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Error().Err(err).Msg("Error trying to read input")
					return err
				}
				eof = true
			}
		}
	}

}
