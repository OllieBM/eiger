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
