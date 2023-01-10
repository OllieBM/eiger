package delta

import (
	"crypto/md5"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {

	source := "Hello"
	target := "HelloWorld"
	sig, err := signature.Calculate(strings.NewReader(source), 5, md5.New())
	require.NoError(t, err)
	require.NotNil(t, sig)

	opW := operation.OpWriter{}
	err = Calculate(strings.NewReader(target), sig, md5.New(), 5, opW)
	require.NoError(t, err)
}
