package operation

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type StringWriteCloser struct {
	strings.Builder
	closed bool
}

func (s *StringWriteCloser) Write(p []byte) (int, error) {
	if s.closed {
		return nil, ErrClosed
	}
	return s.Write(p)
}
func (s *StringWriteCloser) Close() error {
	s.closed = true
	return nil
}
func NewStringWriteCloser() io.WriteCloser {
	return &StringWriteCloser{
		closed: false,
	}
}

func TestNewDiffWriter(t *testing.T) {

	sb := strings.Builder{}
	dw := NewDiffWriter(&sb)
	require.NotNil(t, dw)
	ndw, ok := dw.(*customDiffWriter)
	require.True(t, ok)
	require.NotNil(t, ndw)
}

func TestAddMatch(t *testing.T) {
	sb := strings.Builder{}
	_dw := NewDiffWriter(&sb)
	dw, _ := _dw.(*customDiffWriter)
	dw.AddMatch(0)
	// we don't have anything flush content on first add
	require.Equal(t, sb.Len(), 0)

	require.Equal(t, OpMatch, dw.prevOp.operation)
	require.Equal(t, uint64(0), dw.prevOp.blockIndex)
	require.Nil(t, dw.prevOp.data)
	dw.AddMatch(0)
	// we should have the first block written
	require.NotEqual(t, sb.Len(), 0)
	require.Equal(t, sb.Len(), len("= BLOCK_0\n"))
}

func TestAddMiss(t *testing.T) {
	sb := strings.Builder{}
	_dw := NewDiffWriter(&sb)
	dw, _ := _dw.(*customDiffWriter)
	dw.AddMiss('a')
	// we don't have anything flush content on first command
	require.Equal(t, sb.Len(), 0)
	dw.AddMiss('b')
	// we still don't have anything to flush since we
	// join misses in the writer
	require.Equal(t, sb.Len(), 0)
	dw.AddMiss('c')
	// we still don't have anything to flush since we
	// join misses in the writer
	require.Equal(t, sb.Len(), 0)
	require.Equal(t, OpMiss, dw.prevOp.operation)
	require.Equal(t, uint64(0), dw.prevOp.blockIndex)
	require.Equal(t, []byte("abc"), dw.prevOp.data)

	require.NoError(t, dw.Flush())
	require.Nil(t, dw.prevOp)
	require.NotEqual(t, sb.Len(), 0)
	require.Equal(t, sb.Len(), len("+ 3 abc\n"))
}
