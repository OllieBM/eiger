package operation

import (
	"io"
	"strings"
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
