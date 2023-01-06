package signature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadAsChunks(t *testing.T) {

	in := "123456789"
	chunks := ReadAsChunks(strings.NewReader(in), 1)

	for i, c := range []byte(in) {
		require.Equal(t, byte(c), chunks[i])
	}

}
