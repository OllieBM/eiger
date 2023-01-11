package hasher

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/stretchr/testify/require"
)

// func TestCalculate(t *testing.T) {

// 	data := "this is some input"
// 	in := strings.NewReader(data)
// 	chunkSize := 4
// 	var expect []string

// 	// this chunks data
// 	for len(data) > 0 {
// 		n := chunkSize
// 		if n > len(data) {
// 			n = len(data)
// 		}
// 		expect = append(expect, string(data[:n]))
// 		data = data[n:]
// 	}

// 	hasher := rollingHasher{
// 		reader:    bufio.NewReader(in),
// 		chunkSize: uint32(chunkSize),
// 		buffer:    make([]byte, chunkSize),
// 		rc:        rolling_checksum.New(),
// 	}
// 	for i, exp := range expect {
// 		_, err := hasher.Calculate()
// 		require.Equal(t, exp, string(hasher.buffer), i)
// 		if i != len(expect)-1 {
// 			require.NoError(t, err, i)
// 		} else {
// 			require.Error(t, io.EOF)
// 			//require.Equal(t, uint32(0), h)
// 		}
// 	}

// }

func helperChunkString(in string, chunkSize int) (result []string) {

	// this chunks data
	for len(in) > 0 {
		n := chunkSize
		if n > len(in) {
			n = len(in)
		}
		result = append(result, string(in[:n]))
		in = in[n:]
	}
	return
}

func TestCalculate(t *testing.T) {

	tcs := []struct {
		data      string
		chunkSize uint32
	}{
		{
			"this is some input",
			4,
		},
		{
			"this is some input",
			1,
		},
		{
			"this is some input",
			10,
		},
		{
			"this is some input",
			100,
		},
	}

	for _, tc := range tcs {
		expected := helperChunkString(tc.data, int(tc.chunkSize))

		hasher := rollingHasher{
			reader:    bufio.NewReader(strings.NewReader(tc.data)),
			chunkSize: uint32(tc.chunkSize),
			buffer:    make([]byte, tc.chunkSize),
			rc:        rolling_checksum.New(),
		}

		for _, exp := range expected {
			_, err := hasher.Calculate()
			require.Equal(t, exp, string(hasher.buffer))
			require.NoError(t, err)

		}
		h, err := hasher.Calculate()
		require.ErrorIs(t, err, io.EOF)
		require.Equal(t, uint32(0), h)

	}

}
