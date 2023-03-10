package cmd

import (
	"crypto/md5"
	"os"
	"strings"

	"github.com/OllieBM/eiger/pkg/delta"
	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// https://github.com/madler/zlib/blob/master/adler32.c
// NMAX is the largest n such that 255n(n+1)/2 + (n+1)(BASE-1) <= 2^32-1
// added here for future implementation of rolling adler32
const NMAX = 5552
const defaultBlockSize = 4

var (
	loglevel  string
	output    string // filename to output to
	blockSize uint32
)

func initFlags() {
	// handle flags
	switch strings.ToUpper(loglevel) {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		log.Error().Msgf("Unkown log level %s", loglevel)
	}

	if blockSize > NMAX || blockSize <= 0 {

		log.Error().Msgf("invalid parameter %d for Blocksize using default %d", blockSize, defaultBlockSize)
		blockSize = defaultBlockSize
	}
}

var DiffCmd = &cobra.Command{
	Short: "Create a diff of two files",
	Long:  "Diff File1 against File2 creating a diff file with instructions on how to transform File1 into File2",
	Use:   `diff File1 File2 [flags]`,
	Args:  cobra.MinimumNArgs(2), //expect 2 positional arguments File1 and File2
	RunE: func(cmd *cobra.Command, args []string) error {

		initFlags()

		source, err := os.Open(args[0])
		if err != nil {
			log.Error().Err(err).Msgf("could not open file `%s`", args[0])
			return err
		}
		target, err := os.Open(args[1])
		if err != nil {
			log.Error().Err(err).Msgf("could not open file `%s`", args[1])
			return err
		}

		// create the strong hasher
		// TODO: this can be configurable to different types
		hasher := md5.New()
		rc := rolling_checksum.New()
		sig, err := signature.New(source, int(blockSize), hasher, rc)
		if err != nil {
			return err
		}

		rc = rolling_checksum.New()
		out := os.Stdout
		if output != "" {
			out, err = os.Create(output)
			if err != nil {
				log.Error().Err(err).Msgf("could not open file `%s`", output)
				return err
			}
		}
		//writer := operation.NewDiffWriter(out)
		writer := operation.NewDiffWriter(out)
		err = delta.Calculate(target, sig, rc, writer)
		if err != nil {
			log.Error().Err(err).Msgf("error calculating delta")
			return err
		}
		log.Info().Msg("flushing")
		// err = diffW.Output(os.Stdout)
		err = writer.Flush() // should be something like close
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {

	// tells the DiffCmd to read these flags when it's being used
	DiffCmd.Flags().Uint32VarP(&blockSize, "blocksize", "b", defaultBlockSize, "the size of chunks in bytes to use when matching data from the files max is 0 < b <=5552")
	// 5552 is the maximum value that the rolling checksum algorithm will work for certain algorithms,
	// so we use that as a sensible max
	// /* NMAX is the largest n such that 255n(n+1)/2 + (n+1)(BASE-1) <= 2^32-1 */
	DiffCmd.Flags().StringVarP(&loglevel, "loglevel", "l", "ERROR", "log level to display {DEBUG|INFO|ERROR} default=ERROR")
	DiffCmd.Flags().StringVarP(&output, "output", "o", "", "optional file to write output to")

}

func Execute() error {
	return DiffCmd.Execute()
}
