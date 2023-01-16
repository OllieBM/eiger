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
	sig, err := signature.New(strings.NewReader(old), 5, hasher)
	require.NoError(t, err)
	require.NotNil(t, sig)

	opW = &operation.OpWriter{}
	err = delta.Calculate(strings.NewReader(new), sig, opW)
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
			// add a desc
			"hello",
			"hello",
			"= BLOCK_0\n",
			5,
		},
		{
			// add a desc
			"hello Worl",
			"hello Worl",
			"= BLOCK_0\n= BLOCK_1\n",
			5,
		},
		{
			// add a desc
			// example with tail matching
			"hello World",
			"hello World", // hello| worl |d
			"= BLOCK_0\n= BLOCK_1\n= BLOCK_2\n",
			5,
		},
		{
			// add a desc
			"hello",
			"helloWorld",
			"= BLOCK_0\n+ 5 World\n",
			5,
		},
		{
			// add a desc
			"hello",
			"",
			"",
			5,
		},
		{
			// add a desc

			// this is an example
			// of 'dropping the tail' the check for the ` end`
			// does not match up
			"start end",
			"start middle end", // start| midd|le en|d
			//"= BLOCK_0\n+ 11  middle end\n", // take note of the whitespace ` middle end`
			"= BLOCK_0\n+ 7  middle\n= BLOCK_1\n",
			5,
		},
		{
			// add a desc
			"",
			"123 add missing text 456",
			"+ 24 123 add missing text 456\n",
			5,
		},
	}

	for _, tc := range tcs {

		hasher := md5.New()
		sig, err := signature.New(strings.NewReader(tc.source), 5, hasher)
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
