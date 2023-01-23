package signature

import (
	"crypto/md5"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/stretchr/testify/require"
)

func TestCreateSignature(t *testing.T) {

	in := "HelloWorld"
	reader := strings.NewReader(in)

	s := len(in) / 2
	hasher := md5.New() // TODO: use mockgen tests
	rc := rolling_checksum.New()
	sig, err := New(reader, s, hasher, rc)
	require.NoError(t, err)
	require.Len(t, sig.m, 2) // assumes the weak hashes don't cause collision

	weakHasher := rolling_checksum.New()
	weak := weakHasher.Calculate([]byte(in[0:s]))
	require.Contains(t, sig.m, weak)
	require.Len(t, sig.m[weak], 1)
	weak = weakHasher.Calculate([]byte(in[s:]))
	require.Contains(t, sig.m, weak)
	require.Len(t, sig.m[weak], 1)

}

func TestFindMatch(t *testing.T) {
	source := "Hello"
	hasher := md5.New()
	rc := rolling_checksum.New()
	sig, err := New(strings.NewReader(source), 5, hasher, rc)
	hasher.Reset()
	require.NoError(t, err)
	require.NotNil(t, sig)
	require.Len(t, sig.m, 1)

	weak := rolling_checksum.New().Calculate([]byte(source))
	hasher.Reset()
	found, indx := sig.FindMatch(weak, []byte(source))
	require.True(t, found)
	require.Equal(t, 0, indx)
}
