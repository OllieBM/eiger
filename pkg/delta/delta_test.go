package delta

import (
	"crypto/md5"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/operation"
	"github.com/OllieBM/eiger/pkg/rolling_checksum"
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
func TestFindMatch(t *testing.T) {
	source := "Hello"
	//target := "HelloWorld"
	hasher := md5.New()
	sig, err := signature.Calculate(strings.NewReader(source), 5, hasher)
	hasher.Reset()
	require.NoError(t, err)
	require.NotNil(t, sig)
	require.Len(t, sig, 1)

	weak := rolling_checksum.New().Calculate([]byte(source))
	hasher.Reset()
	found, indx := FindMatch(weak, []byte(source), hasher, sig)
	require.True(t, found)
	require.Equal(t, 0, indx)
}

func TestCalculate(t *testing.T) {

	source := "Hello"
	target := "HelloWorld"
	sig, err := signature.Calculate(strings.NewReader(source), 5, md5.New())
	require.NoError(t, err)
	require.NotNil(t, sig)

	opW := operation.OpWriter{} // use mockgen
	err = Calculate(strings.NewReader(target), sig, md5.New(), 5, &opW)
	require.NoError(t, err)
	require.Len(t, opW.Operations(), 2)

}

func TestE2E(t *testing.T) {
	old := "Hello"
	new := "HelloWorld"

	opW := &operation.OpWriter{}
	opW.AddMatch(0)
	expected := opW.Operations()
	require.Len(t, expected, 1)
	opW.AddMiss('W')
	opW.AddMiss('o')
	opW.AddMiss('r')
	opW.AddMiss('l')
	opW.AddMiss('d')
	expected = opW.Operations()
	require.Len(t, expected, 2)
	// create a signature

	hasher := md5.New()
	sig, err := signature.Calculate(strings.NewReader(old), 5, hasher)
	require.NoError(t, err)
	require.NotNil(t, sig)

	opW = &operation.OpWriter{}
	err = Calculate(strings.NewReader(new), sig, hasher, 5, opW)
	require.NoError(t, err)

	require.Equal(t, expected, opW.Operations())

}
