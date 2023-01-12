package main

import (
	"crypto/md5"
	"os"

	"github.com/OllieBM/eiger/pkg/delta"
	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// try to open file one
// try to open file two
// create signature from file one
// create delta for file two

func main() {
	var file1, file2 string
	var blockSize uint32
	// does not use the general cobra layout,
	// but still makes passing in values easier
	var rootCmd = &cobra.Command{
		Short: "Create a diff of two files",
		Use: `Diff File1 against File2 creating a diff file with instructions on how to transform File1 into File2
		eiger-diff --file1 <File1> --file2 <File2>`, //eiger-diff --file1 <File1> --file2 <File2>``,
		//Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			f1, err := os.Open(file1)
			if err != nil {
				log.Error().Err(err)
			}
			f2, err := os.Open(file2)
			if err != nil {
				log.Error().Err(err)
			}

			// create the strong hasher
			hasher := md5.New()

			sig, err := signature.Calculate(f2, int(blockSize), hasher)
			if err != nil {
				log.Error().Err(err)
				return err
			}
			opW := &operation.OpWriter{}
			err = delta.Calculate(f1, sig, hasher, uint64(blockSize), opW)
			if err != nil {
				log.Error().Err(err)
				return err
			}
			err = opW.Flush(os.Stdout)
			if err != nil {
				log.Error().Err(err)
				return err
			}
			return nil
		},
	}
	// add flags to the command
	rootCmd.Flags().StringVarP(&file1, "file1", "", "", "The first file to read")
	rootCmd.Flags().StringVarP(&file2, "file2", "", "", "The second file to read")
	rootCmd.MarkFlagRequired("file1") // if not supplied will panic
	rootCmd.MarkFlagRequired("file2") // if not supplied will panic
	rootCmd.Flags().Uint32VarP(&blockSize, "blocksize", "b", 32, "the size of chunks in bytes to use when matching data from the files")
	// TODO: add in stream output
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}
	os.Exit(0)
}
