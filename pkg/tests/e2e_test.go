package e2e_test

import (
	"crypto/md5"
	"io"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/delta"
	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/stretchr/testify/require"
)

// /////// E2E TESTS
// func TestE2E(t *testing.T) {
// 	old := "Hello"
// 	new := "HelloWorld"

// 	opW := &operation.OpWriter{}
// 	opW.AddMatch(0)
// 	expected := opW.Operations()
// 	require.Len(t, expected, 1)
// 	opW.AddMiss('W')
// 	opW.AddMiss('o')
// 	opW.AddMiss('r')
// 	opW.AddMiss('l')
// 	opW.AddMiss('d')
// 	expected = opW.Operations()
// 	require.Len(t, expected, 2)
// 	// create a signature

// 	rc := rolling_checksum.New()
// 	hasher := md5.New()
// 	sig, err := signature.New(strings.NewReader(old), 5, rc, hasher)
// 	require.NoError(t, err)
// 	require.NotNil(t, sig)

// 	opW = &operation.OpWriter{}
// 	//err = delta.Calculate(strings.NewReader(new), sig, rc, opW)
// 	require.NoError(t, err)

// 	require.Equal(t, expected, opW.Operations())

// 	//opW.Output()
// }

// func TestOutput(t *testing.T) {

// 	tcs := []struct {
// 		description string
// 		source      string
// 		target      string
// 		expected    string
// 		chunkSize   int
// 	}{
// 		{
// 			"matching_files_of_blocksize_length",
// 			"hello",
// 			"hello",
// 			"= BLOCK_0\n",
// 			5,
// 		},
// 		{
// 			"Matching_files_of_multiple_blocksize",
// 			"hello Worl",
// 			"hello Worl",
// 			"= BLOCK_0\n= BLOCK_1\n",
// 			5,
// 		},
// 		{
// 			"matching_files_with_tail_<_blocksize",
// 			"hello World",
// 			"hello World", // hello| worl |d
// 			"= BLOCK_0\n= BLOCK_1\n= BLOCK_2\n",
// 			5,
// 		},
// 		{
// 			"rearranged_blocks",
// 			"helloWorld",
// 			"Worldhello",
// 			"= BLOCK_1\n= BLOCK_0\n", // BLOCK_2 is '\n'
// 			5,
// 		},
// 		{
// 			"rearranged_blocks_with_tail_addition",
// 			"helloWorld",
// 			"WorldhelloMore",
// 			"= BLOCK_1\n= BLOCK_0\n+ 4 More\n", // BLOCK_2 is '\n'
// 			5,
// 		},
// 		{
// 			"rearranged_blocks_with_middle_addition",
// 			"helloWorld",
// 			"WorldMorehello",
// 			"= BLOCK_1\n+ 4 More\n= BLOCK_0\n", // BLOCK_2 is '\n'
// 			5,
// 		},
// 		{
// 			"rearranged blocks_with_start_addition",
// 			"helloWorld",
// 			"MoreWorldhello",
// 			"+ 4 More\n= BLOCK_1\n= BLOCK_0\n", // BLOCK_2 is '\n'
// 			5,
// 		},
// 		{
// 			"rearranged_blocks_with_removal",
// 			"helloMoreWorld", // |hello|MoreW|orld since "World will not be matched since it is not a tail"
// 			"Worldhello",
// 			"+ 5 World\n= BLOCK_0\n",
// 			5,
// 		},
// 		{
// 			"rearranged_blocks_with_removal_uniform_blocks",
// 			"helloMore World", // |hello|More |World
// 			"Worldhello",
// 			"= BLOCK_2\n= BLOCK_0\n",
// 			5,
// 		},
// 		{
// 			"rearranged_blocks_with_removal",
// 			"helloMoreWorld", // |hello|MoreW|orld
// 			"Worldhello",
// 			"+ 5 World\n= BLOCK_0\n",
// 			5,
// 		},
// 		{
// 			"additional_characters",
// 			"hello",
// 			"helloWorld",
// 			"= BLOCK_0\n+ 5 World\n",
// 			5,
// 		},
// 		{
// 			"removed_charactesrs",
// 			"hello",
// 			"",
// 			"", // empty delta, because we don't want any references or addition blocks
// 			5,
// 		},
// 		{
// 			"removed_charactesrs",
// 			"helloworld",
// 			"world",
// 			"= BLOCK_1\n", // empty delta, because we don't want any references or additional
// 			5,
// 		},
// 		{
// 			// this is an example
// 			// of 'dropping the tail' the check for the ` end`
// 			// does not match up
// 			"addition_between_blocks",
// 			"start end",
// 			"start middle end", // start| midd|le en|d => matched as |start| middle| end

// 			"= BLOCK_0\n+ 7  middle\n= BLOCK_1\n",
// 			5,
// 		},
// 		{
// 			"empty_source_file_(all_additions)",
// 			"",
// 			"123 add missing text 456",
// 			"+ 24 123 add missing text 456\n",
// 			5,
// 		},
// 	}

// 	for _, tc := range tcs {

// 		t.Run(tc.description, func(t *testing.T) {
// 			hasher := md5.New()
// 			rc := rolling_checksum.New()
// 			sig, err := signature.New(strings.NewReader(tc.source), 5, hasher, rc)
// 			require.NoError(t, err)
// 			require.NotNil(t, sig)

// 			opW := &operation.OpWriter{}
// 			err = delta.Calculate2(strings.NewReader(tc.target), sig, opW)
// 			require.NoError(t, err)

// 			out := &strings.Builder{}
// 			err = opW.Output(out)
// 			require.NoError(t, err)
// 			require.Equal(t, tc.expected, out.String())
// 		})
// 	}

// }

///////////////////////////

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
	rc := rolling_checksum.New()

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
			// this would count as a missing start, then an added start at the index of the next block that would be read by the applier
			"moved_chunk_start_to_end",
			"Hello12345World",
			"12345WorldHello",
			"- BLOCK_0\n= @3 BLOCK_0\n",
			5,
		},
		{
			// here two blocks have been swapped, but its registerd
			// as a sequencial remove, then additions
			"moved_chunk_mid_to_mid",
			"Hello01234abcdeWorld56789",
			"HelloWorldabcde0123456789",
			"- BLOCK_1\n- BLOCK_2\n= @4 BLOCK_2\n= @4 BLOCK_1\n",

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

		{
			"added_chars_start",
			"WorldItsme",
			"HelloWorldItsme",
			"+ @0 5 Hello\n", // for of + @<offset of the next block to read>  addition
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
			// this would mean when applying you read blocks (pos -> @N) [0,1] add the (5) characters 'Itsme' and then continue reading blocks 2:END
			// fast forward to block 2 start, add 'Itsme'
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
			"+ @0 5 Hello\n+ @2 5 Again\n", // add 'Hello" before block 0, .. Before block 2 add Again, (BLOCK 2 doesn't exist so appliert would end exit)
			5,
		},
		{
			"added_chars_middle_end",
			"HelloItsme",
			"HelloWorldItsmeAgain",
			"+ @1 5 World\n+ @2 5 Again\n", // before block 1 add 'World', (fast forward from block 1...block 2 start) before block 2 add 'Again'
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
		{
			// if we make the assumption a delta applier will copy all unreferenced blocks, then any removals will need to be transmitted
			// only the fist block
			"trimmed_file",
			`0123456789abcdefghijklmnopqrst`,
			`0123456789`,
			"- BLOCK_1\n- BLOCK_2\n",
			10,
		},
		{
			// if we make the assumption a delta applier will copy all unreferenced blocks, then any removals will need to be transmitted
			// this example the source will be split into
			"trimmed_file_repeating_text",
			`123456789123456789123456789123456789`,
			`123456789123456789`,
			"= @1 BLOCK_0\n- BLOCK_1\n- BLOCK_2\n- BLOCK_3\n", // this case can be problematic since it depends how the hasher matches the file
			9, // a lower threshold can cause collisions in the weak_sum and strong sum, consider looking into a checksum which also uses the block_offset
		},
		{
			"repeating_text_larger_example",
			`12346789
12346789
12346789
12346789
123456789
123456789
123456789
123456789`,
			`123456789
123456789
123456789
123456789
123456789
123456789
123456789
123456789`,
			"- BLOCK_0\n- BLOCK_1\n- BLOCK_2\n- BLOCK_3\n+ @8 1 7\n= @8 BLOCK_7\n+ @8 1 7\n= @8 BLOCK_7\n+ @8 1 7\n= @8 BLOCK_7\n+ @8 1 7\n= @8 BLOCK_7\n",
			9,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			sb := strings.Builder{}
			s := NewStringWriteCloser(&sb)
			writer := operation.NewMinimalDiffWriter(s)

			sig, err := signature.New(strings.NewReader(tc.source), tc.chunkSize, hasher, rc)
			require.NoError(t, err)
			require.NotNil(t, sig)

			err = delta.Calculate3(strings.NewReader(tc.target), sig, rc, writer)
			require.NoError(t, err)
			require.NoError(t, writer.Flush())
			require.Equal(t,
				tc.expected,
				sb.String(),
			)

		})
	}

}
