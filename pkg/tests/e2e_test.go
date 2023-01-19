package e2e_test

import (
	"crypto/md5"
	"io"
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
		description string
		source      string
		target      string
		expected    string
		chunkSize   int
	}{
		{
			"matching_files_of_blocksize_length",
			"hello",
			"hello",
			"= BLOCK_0\n",
			5,
		},
		{
			"Matching_files_of_multiple_blocksize",
			"hello Worl",
			"hello Worl",
			"= BLOCK_0\n= BLOCK_1\n",
			5,
		},
		{
			"matching_files_with_tail_<_blocksize",
			"hello World",
			"hello World", // hello| worl |d
			"= BLOCK_0\n= BLOCK_1\n= BLOCK_2\n",
			5,
		},
		{
			"rearranged_blocks",
			"helloWorld",
			"Worldhello",
			"= BLOCK_1\n= BLOCK_0\n", // BLOCK_2 is '\n'
			5,
		},
		{
			"rearranged_blocks_with_tail_addition",
			"helloWorld",
			"WorldhelloMore",
			"= BLOCK_1\n= BLOCK_0\n+ 4 More\n", // BLOCK_2 is '\n'
			5,
		},
		{
			"rearranged_blocks_with_middle_addition",
			"helloWorld",
			"WorldMorehello",
			"= BLOCK_1\n+ 4 More\n= BLOCK_0\n", // BLOCK_2 is '\n'
			5,
		},
		{
			"rearranged blocks_with_start_addition",
			"helloWorld",
			"MoreWorldhello",
			"+ 4 More\n= BLOCK_1\n= BLOCK_0\n", // BLOCK_2 is '\n'
			5,
		},
		{
			"rearranged_blocks_with_removal",
			"helloMoreWorld", // |hello|MoreW|orld since "World will not be matched since it is not a tail"
			"Worldhello",
			"+ 5 World\n= BLOCK_0\n",
			5,
		},
		{
			"rearranged_blocks_with_removal_uniform_blocks",
			"helloMore World", // |hello|More |World
			"Worldhello",
			"= BLOCK_2\n= BLOCK_0\n",
			5,
		},
		{
			"rearranged_blocks_with_removal",
			"helloMoreWorld", // |hello|MoreW|orld
			"Worldhello",
			"+ 5 World\n= BLOCK_0\n",
			5,
		},
		{
			"additional_characters",
			"hello",
			"helloWorld",
			"= BLOCK_0\n+ 5 World\n",
			5,
		},
		{
			"removed_charactesrs",
			"hello",
			"",
			"", // empty delta, because we don't want any references or addition blocks
			5,
		},
		{
			"removed_charactesrs",
			"helloworld",
			"world",
			"= BLOCK_1\n", // empty delta, because we don't want any references or additional
			5,
		},
		{
			// this is an example
			// of 'dropping the tail' the check for the ` end`
			// does not match up
			"addition_between_blocks",
			"start end",
			"start middle end", // start| midd|le en|d => matched as |start| middle| end

			"= BLOCK_0\n+ 7  middle\n= BLOCK_1\n",
			5,
		},
		{
			"empty_source_file_(all_additions)",
			"",
			"123 add missing text 456",
			"+ 24 123 add missing text 456\n",
			5,
		},
	}

	for _, tc := range tcs {

		t.Run(tc.description, func(t *testing.T) {
			hasher := md5.New()
			sig, err := signature.New(strings.NewReader(tc.source), 5, hasher)
			require.NoError(t, err)
			require.NotNil(t, sig)

			opW := &operation.OpWriter{}
			err = delta.Calculate(strings.NewReader(tc.target), sig, opW)
			require.NoError(t, err)

			out := &strings.Builder{}
			err = opW.Output(out)
			require.NoError(t, err)
			require.Equal(t, tc.expected, out.String())
		})
	}

}

/////

type StringWriteCloser struct {
	*strings.Builder
	closed bool
}

func (s *StringWriteCloser) Write(p []byte) (int, error) {
	if s.closed {
		return 0, operation.ErrClosed
	}
	return s.Builder.Write(p)
}
func (s *StringWriteCloser) Close() error {
	s.closed = true
	return nil
}
func NewStringWriteCloser(builder *strings.Builder) io.WriteCloser {
	return &StringWriteCloser{
		Builder: builder,
		closed:  false,
	}
}

func TestOnlyDiffE2E(t *testing.T) {

	hasher := md5.New()

	tcs := []struct {
		desc      string
		source    string
		target    string
		expected  string
		chunkSize int
	}{
		{
			"matching_files",
			"helloworld",
			"helloworld",
			"",
			5,
		},
		{
			"matching_files_with_tail",
			"helloworldtail", // tail is < blocksize
			"helloworldtail",
			"", // tail matching is not ideal
			5,
		},
		{
			"empty_target",
			"helloworld",
			"",
			"- BLOCK_0\n- BLOCK_1\n",
			5,
		},
		{
			"empty_target_with_tail",
			"helloworldtail",
			"",
			"- BLOCK_0\n- BLOCK_1\n- BLOCK_2\n",
			5,
		},

		{
			"missing_chunk_start",
			"HelloWorldItsme",
			"WorldItsme",
			"- BLOCK_0\n",
			5,
		},
		{
			"missing_chunk_middle",
			"HelloWorldItsme",
			"HelloItsme",
			"- BLOCK_1\n",
			5,
		},
		{
			"missing_chunk_end",
			"HelloWorldItsme",
			"HelloWorld",
			"- BLOCK_2\n",
			5,
		},
		// moved chunk
		{
			// this would count as a missing start, then an added start
			"moved_chunk_start_to_end",
			"Hello12345World",
			"12345WorldHello",
			"- BLOCK_0\n= @3 BLOCK_0\n",
			5,
		},
		{
			// counts as missing block0 & 1  then reference block 0 & 1
			// offset is the block index to add things to (before that specific block)
			// this example @3 would mean, add/reference block before reading block 3
			"moved_chunk_end_to_start",
			"HelloWorld12345",
			"12345HelloWorld",
			"- BLOCK_0\n- BLOCK_1\n= @3 BLOCK_0\n= @3 BLOCK_1\n",
			5,
		},

		// TODO: fix
		{
			"added_chars_start",
			"WorldItsme",
			"HelloWorldItsme",
			"+ @0 5 Hello\n",
			5,
		},
		{
			"added_chars_middle",
			"HelloItsme",
			"HelloWorldItsme",
			"+ @1 5 World\n",
			5,
		},
		{
			"added_chars_end",
			"HelloWorld",
			"HelloWorldItsme",
			"+ @2 5 Itsme\n",
			5,
		},
		// multiple added
		{
			"added_chars_start_and_end",
			"WorldItsme",
			"HelloWorldItsmeAgain",
			"+ @0 5 Hello\n+ @2 5 Again\n",
			5,
		},
		{
			"added_chars_middle_end",
			"HelloItsme",
			"HelloWorldItsmeAgain",
			"+ @1 5 World\n+ @2 5 Again\n",
			5,
		},
		// trailing removes
		{
			"trailing_removed",
			"HelloWorldItsMe12345",
			"HelloWorld",
			"- BLOCK_2\n- BLOCK_3\n",
			5,
		},
		{
			"trailing_removed_with_removed_start",
			"HelloWorldItsMe12345",
			"World",
			"- BLOCK_0\n- BLOCK_2\n- BLOCK_3\n",
			5,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			sb := strings.Builder{}
			s := NewStringWriteCloser(&sb)
			writer := operation.NewOnlyDiffWriter(s)

			sig, err := signature.New(strings.NewReader(tc.source), 5, hasher)
			require.NoError(t, err)
			require.NotNil(t, sig)

			err = delta.Calculate3(strings.NewReader(tc.target), sig, writer)
			require.NoError(t, err)
			require.NoError(t, writer.Flush())
			require.Equal(t,
				tc.expected,
				sb.String(),
			)

		})
	}

}
