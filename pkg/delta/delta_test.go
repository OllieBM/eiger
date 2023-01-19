package delta

import (
	"crypto/md5"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/stretchr/testify/require"
)

func TestMD5(t *testing.T) {
	hasher := md5.New()
	a := hasher.Sum([]byte("Hello"))
	hasher.Reset()
	b := hasher.Sum([]byte("Hello"))
	require.Equal(t, a, b)

}

func TestCalculate(t *testing.T) {

	source := "Hello"
	target := "HelloWorld"
	sig, err := signature.New(strings.NewReader(source), 5, md5.New())
	require.NoError(t, err)
	require.NotNil(t, sig)

	opW := operation.OpWriter{} // use mockgen
	err = Calculate(strings.NewReader(target), sig, &opW)
	require.NoError(t, err)
	require.Len(t, opW.Operations(), 2)

}

func TestCalculate2(t *testing.T) {

	source := "Hello"
	target := "HelloWorld"
	sig, err := signature.New(strings.NewReader(source), 5, md5.New())
	require.NoError(t, err)
	require.NotNil(t, sig)

	opW := operation.OpWriter{} // use mockgen
	err = Calculate2(strings.NewReader(target), sig, &opW)
	require.NoError(t, err)
	//require.ErrorIs(t, err, io.EOF)
	require.Len(t, opW.Operations(), 2)
}
