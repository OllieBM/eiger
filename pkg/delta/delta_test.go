package delta

import (
	"crypto/md5"
	"io"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/mocks"
	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	source := "Hello"
	target := "HelloWorld"

	rc := rolling_checksum.New()
	sig, err := signature.New(strings.NewReader(source), 5, md5.New(), rc)
	require.NoError(t, err)
	require.NotNil(t, sig)

	mockWriter := mocks.NewMockDiffWriter(ctrl)
	mockWriter.EXPECT().AddMatch(uint64(0)).Times(1)
	mockWriter.EXPECT().AddMiss(byte('W')).Times(1)
	mockWriter.EXPECT().AddMiss(byte('o')).Times(1)
	mockWriter.EXPECT().AddMiss(byte('r')).Times(1)
	mockWriter.EXPECT().AddMiss(byte('l')).Times(1)
	mockWriter.EXPECT().AddMiss(byte('d')).Times(1)
	mockWriter.EXPECT().AddMissingIndexes([]uint64{}).Times(1)
	err = Calculate(strings.NewReader(target), sig, rc, mockWriter)
	require.NoError(t, err)
}

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

func TestDeltaE2E(t *testing.T) {

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
			"- CHUNK_0...CHUNK_2\n",
			5,
		},
		{
			"empty_target_with_tail",
			"helloworldtail",
			"",
			"- CHUNK_0...CHUNK_3\n",
			5,
		},
		{
			"missing_chunk_start",
			"HelloWorldItsme",
			"WorldItsme",
			"- CHUNK_0...CHUNK_1\n",
			5,
		},
		{
			"missing_chunk_middle",
			"HelloWorldItsme",
			"HelloItsme",
			"- CHUNK_1...CHUNK_2\n",
			5,
		},
		{
			"missing_chunk_end",
			"HelloWorldItsme",
			"HelloWorld",
			"- CHUNK_2...CHUNK_3\n",
			5,
		},
		// moved chunk
		{
			// this would count as a missing start, then an added start at the index of the next block that would be read by the applier
			"moved_chunk_start_to_end",
			"Hello12345World",
			"12345WorldHello",
			"- CHUNK_0...CHUNK_1\n= @3 CHUNK_0\n",
			5,
		},
		{
			// here two blocks have been swapped, but its registerd
			// as a sequencial remove, then additions
			"moved_chunk_mid_to_mid",
			"Hello01234abcdeWorld56789",
			"HelloWorldabcde0123456789",
			"- CHUNK_1...CHUNK_3\n= @4 CHUNK_2\n= @4 CHUNK_1\n",

			5,
		},
		{
			// counts as missing block0 & 1  then reference block 0 & 1
			// offset is the block index to add things to (before that specific block)
			// this example @3 would mean, add/reference block before reading block 3
			"moved_chunk_end_to_start",
			"HelloWorld12345",
			"12345HelloWorld",
			"- CHUNK_0...CHUNK_2\n= @3 CHUNK_0\n= @3 CHUNK_1\n",
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
			"- CHUNK_2...CHUNK_4\n",
			5,
		},
		{
			"trailing_removed_with_removed_start",
			"HelloWorldItsMe12345",
			"World",
			"- CHUNK_0...CHUNK_1\n- CHUNK_2...CHUNK_4\n",
			5,
		},
		{
			// if we make the assumption a delta applier will copy all unreferenced blocks, then any removals will need to be transmitted
			// only the fist block
			"trimmed_file",
			`0123456789abcdefghijklmnopqrst`,
			`0123456789`,
			"- CHUNK_1...CHUNK_3\n",
			10,
		},
		{
			// if we make the assumption a delta applier will copy all unreferenced blocks, then any removals will need to be transmitted
			// this example the source will be split into  3 identical blocks
			// this case can be problematic since it depends how the hasher matches the file
			"trimmed_file_repeating_text",
			`123456789123456789123456789123456789`,
			`123456789123456789`,
			"= @1 CHUNK_0\n- CHUNK_1...CHUNK_4\n",
			9, // a lower threshold can cause collisions in the weak_sum and strong sum, consider looking into a checksum which also uses the CHUNK_offset
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
			"- CHUNK_0...CHUNK_4\n+ @8 1 7\n= @8 CHUNK_7\n+ @8 1 7\n= @8 CHUNK_7\n+ @8 1 7\n= @8 CHUNK_7\n+ @8 1 7\n= @8 CHUNK_7\n",
			9,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			sb := strings.Builder{}
			s := NewStringWriteCloser(&sb)
			writer := operation.NewDiffWriter(s)

			sig, err := signature.New(strings.NewReader(tc.source), tc.chunkSize, hasher, rc)
			require.NoError(t, err)
			require.NotNil(t, sig)

			err = Calculate(strings.NewReader(tc.target), sig, rc, writer)
			require.NoError(t, err)
			require.NoError(t, writer.Flush())
			require.Equal(t,
				tc.expected,
				sb.String(),
			)

		})
	}

}
