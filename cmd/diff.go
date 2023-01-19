package cmd

import (
	"crypto/md5"
	"os"
	"strings"

	"github.com/OllieBM/eiger/pkg/delta"
	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const NMAX = 5552

var (
	loglevel  string
	blockSize uint32
)

var DiffCmd = &cobra.Command{
	Short: "Create a diff of two files",
	Long:  "create a diff which can be used to convert a source file (File1) into a target file (File2)",
	Use:   `diff File1 File2 [flags]`,
	Args:  cobra.MinimumNArgs(2), //expect 2 positional arguments File1 and File2
	RunE: func(cmd *cobra.Command, args []string) error {

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
		sig, err := signature.New(source, int(blockSize), hasher)
		if err != nil {
			return err
		}
		opW := &operation.OpWriter{}
		err = delta.Calculate2(target, sig, opW)
		if err != nil {
			log.Error().Err(err).Msgf("error calculating delta")
			return err
		}
		log.Info().Msg("flushing")
		err = opW.Flush(os.Stdout)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {

	// add flags to the command
	//DiffCmd.Flags().StringVarP(&file1, "file1", "", "", "The source file to read")
	//DiffCmd.Flags().StringVarP(&file2, "file2", "", "", "The target file to read")
	// TODO: add in option to write to a file
	//DiffCmd.MarkFlagRequired("file1") // if not supplied will panic
	//DiffCmd.MarkFlagRequired("file2") // if not supplied will panic
	DiffCmd.Flags().Uint32VarP(&blockSize, "blocksize", "b", 5, "the size of chunks in bytes to use when matching data from the files max is 0 < b <=5552")
	// 5552 is the maximum value that the rolling checksum algorithm will work for
	DiffCmd.Flags().StringVarP(&loglevel, "loglevel", "l", "ERROR", "log level to display {DEBUG|INFO|ERROR} default=ERROR")

}

func Execute() {

	log.Logger = log.With().Caller().Logger()
	// left out some log levels
	switch strings.ToUpper(loglevel) {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
	if blockSize > NMAX || 0 < blockSize {
		log.Error().Msg("invalid parameter for Blocksize")
	}

	if err := DiffCmd.Execute(); err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}
}
