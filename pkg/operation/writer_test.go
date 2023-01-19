package operation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOperationFunctional(t *testing.T) {

	opW := OpWriter{}
	opW.AddMiss(byte('a'))
	opW.AddMiss(byte('b'))
	opW.AddMiss(byte('c'))
	require.Len(t, opW.ops, 1)
	opW.AddMatch(0)
	require.Len(t, opW.ops, 2)
	ops := opW.Operations()
	require.Equal(t, OpMiss, ops[0].operation)
	require.Equal(t, ops[0].data, []byte("abc"))
	require.Equal(t, OpMatch, ops[1].operation)
	require.Equal(t, uint64(0), ops[1].blockIndex)

	require.NotEmpty(t, opW.Operations())
	require.Equal(t, opW.ops, opW.Operations())

	expectdDelta := "+ 3 abc\n= BLOCK_0\n"
	sbuilder := strings.Builder{}
	require.NoError(t, opW.Output(&sbuilder))
	require.Equal(t, expectdDelta, sbuilder.String())

}
