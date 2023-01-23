package operation

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// types to help test with
type StringWriteCloser struct {
	*strings.Builder
	closed bool
}

func (s *StringWriteCloser) Write(p []byte) (int, error) {
	if s.closed {
		return 0, ErrClosed
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

func TestNewDiffWriter(t *testing.T) {

	s := NewStringWriteCloser(&strings.Builder{})
	w := NewDiffWriter(s)
	require.NotNil(t, w)
	_, ok := w.(*diffWriter)
	require.True(t, ok)

}

func TestNewDiffWriterAddMatch(t *testing.T) {

	sb := strings.Builder{}
	s := NewStringWriteCloser(&sb)
	w, ok := NewDiffWriter(s).(*diffWriter)
	require.True(t, ok)

	w.AddMatch(0)
	w.Flush()
	require.Equal(t, "", sb.String())
	require.Equal(t, uint64(1), w.blockOffset)
	w.AddMatch(1)
	w.Flush()
	require.Equal(t, "", sb.String())
	require.Equal(t, uint64(2), w.blockOffset)
	// if block 4 was matched before block 2
	// should write out two removals
	w.AddMatch(4)
	w.Flush()
	require.Equal(t, "- CHUNK_2...CHUNK_4\n", sb.String())
	require.Equal(t, uint64(5), w.blockOffset)

	// if block 3 was matched after reading block 4
	w.AddMatch(3)
	w.Flush()
	require.Equal(t, "- CHUNK_2...CHUNK_4\n= @5 CHUNK_3\n", sb.String())
	require.Equal(t, uint64(5), w.blockOffset)
	// if block 4 was matched after reading block 4 & block 4
	w.AddMatch(4)
	w.Flush()

	require.Equal(t, "- CHUNK_2...CHUNK_4\n= @5 CHUNK_3\n= @5 CHUNK_4\n", sb.String())
	require.Equal(t, uint64(5), w.blockOffset)

}

func TestAddMissingIndexes(t *testing.T) {

	t.Run("only_missing", func(t *testing.T) {
		sb := strings.Builder{}
		s := NewStringWriteCloser(&sb)
		w, ok := NewDiffWriter(s).(*diffWriter)
		require.True(t, ok)
		w.AddMissingIndexes([]uint64{0, 1, 2, 3, 4})
		require.Equal(t, "", sb.String())
		require.NoError(t, w.flush())
		require.Equal(t, "- CHUNK_0...CHUNK_5\n", sb.String())
	})

	t.Run("missing_start", func(t *testing.T) {
		sb := strings.Builder{}
		s := NewStringWriteCloser(&sb)
		w, ok := NewDiffWriter(s).(*diffWriter)
		require.True(t, ok)
		w.AddMatch(3)
		w.AddMatch(4)
		require.NoError(t, w.flush())
		require.Equal(t, "- CHUNK_0...CHUNK_3\n", sb.String())
		sb.Reset()

		// these values are less than or equal to the last written
		// block index, so are considered to have been written as missed
		// already
		w.AddMissingIndexes([]uint64{0, 1, 2, 3, 4})
		require.Equal(t, "", sb.String())
		require.NoError(t, w.flush())
		require.Equal(t, "", sb.String())

		// if there were missing from the tail
		w.AddMissingIndexes([]uint64{5, 6, 7})
		require.Equal(t, "", sb.String())
		require.NoError(t, w.flush())
		require.Equal(t, "- CHUNK_5...CHUNK_8\n", sb.String())
	})
}

func TestAddMiss(t *testing.T) {

	t.Run("consecutive", func(t *testing.T) {
		sb := strings.Builder{}
		s := NewStringWriteCloser(&sb)
		w, ok := NewDiffWriter(s).(*diffWriter)
		require.True(t, ok)
		w.AddMiss('a')
		w.AddMiss('b')
		w.AddMiss('c')
		require.Equal(t, "", sb.String())
		w.flush()
		require.Equal(t, "+ @0 3 abc\n", sb.String())
	})
	t.Run("interleaved", func(t *testing.T) {
		sb := strings.Builder{}
		s := NewStringWriteCloser(&sb)
		w, ok := NewDiffWriter(s).(*diffWriter)
		require.True(t, ok)
		w.AddMiss('a')
		w.AddMatch(0)
		require.Equal(t, "+ @0 1 a\n", sb.String())
		w.AddMiss('b')
		require.Equal(t, "+ @0 1 a\n", sb.String())
		w.AddMatch(2)
		require.Equal(t, "+ @0 1 a\n+ @1 1 b\n", sb.String())
		w.AddMiss('c')
		require.Equal(t, "+ @0 1 a\n+ @1 1 b\n- CHUNK_1...CHUNK_2\n", sb.String())
		require.NoError(t, w.Flush())
		require.Equal(t, "+ @0 1 a\n+ @1 1 b\n- CHUNK_1...CHUNK_2\n+ @3 1 c\n", sb.String())

	})
}
