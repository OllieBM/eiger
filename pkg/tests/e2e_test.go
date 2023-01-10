package e2e_test

import (
	"crypto/md5"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/delta"
	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/stretchr/testify/require"
)

/////// E2E TESTS
func TestE2E(t *testing.T) {
	old := "Hello"
	new := "HelloWorld"

	opW := &operation.OpWriter{}
	opW.AddMatch(0)
	expected := opW.Operations()
	require.Len(t, expected, 1)
	opW.AddMiss('W')
	opW.AddMiss('o')
	opW.AddMiss('r')
	opW.AddMiss('l')
	opW.AddMiss('d')
	expected = opW.Operations()
	require.Len(t, expected, 2)
	// create a signature

	hasher := md5.New()
	sig, err := signature.Calculate(strings.NewReader(old), 5, hasher)
	require.NoError(t, err)
	require.NotNil(t, sig)

	opW = &operation.OpWriter{}
	err = delta.Calculate(strings.NewReader(new), sig, hasher, 5, opW)
	require.NoError(t, err)

	require.Equal(t, expected, opW.Operations())

	//opW.Output()
}

func TestOutput(t *testing.T) {

	tcs := []struct {
		source    string
		target    string
		expected  string
		chunkSize int
	}{
		{
			"hello",
			"hello",
			"= BLOCK_0\n",
			5,
		},
		{
			"hello",
			"helloWorld",
			"= BLOCK_0\n+ 5 World\n",
			5,
		},
		{
			"hello",
			"",
			"",
			5,
		},
		{
			// this is an example
			// of 'dropping the tail'
			// we could still maintain using the last values now we support
			// weak hash collisions
			// but current behaviour here is an example of
			"start end",
			"start middle end",
			"= BLOCK_0\n+ 11  middle end\n", // take note of the whitespace ` middle end`
			5,
		},
		{
			"",
			"123 add missing text 456",
			"+ 24 123 add missing text 456\n",
			5,
		},
	}

	for _, tc := range tcs {

		hasher := md5.New()
		sig, err := signature.Calculate(strings.NewReader(tc.source), 5, hasher)
		require.NoError(t, err)
		require.NotNil(t, sig)

		opW := &operation.OpWriter{}
		err = delta.Calculate(strings.NewReader(tc.target), sig, hasher, 5, opW)
		require.NoError(t, err)

		out := &strings.Builder{}
		err = opW.Output(out)
		require.NoError(t, err)
		require.Equal(t, tc.expected, out.String())

	}

}

/////
