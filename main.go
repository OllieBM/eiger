package main

import (
	"os"

	"github.com/OllieBM/eiger/cmd"
	"github.com/rs/zerolog/log"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(-1)
	}
	os.Exit(0)
}
