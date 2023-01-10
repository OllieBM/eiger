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

	sig, err := Calculate(reader, s, hasher)
	require.NoError(t, err)
	require.Len(t, sig, 2) // assumes the weak hashes don't cause collision

	weakHasher := rolling_checksum.New()
	weak := weakHasher.Calculate([]byte(in[0:s]))
	require.Contains(t, sig, weak)
	require.Len(t, sig[weak], 1)
	weak = weakHasher.Calculate([]byte(in[s:]))
	require.Contains(t, sig, weak)
	require.Len(t, sig[weak], 1)

}
