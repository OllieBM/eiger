package operation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOnlyDiffWriter(t *testing.T) {

	s := NewStringWriteCloser(&strings.Builder{})
	w := NewOnlyDiffWriter(s)
	require.NotNil(t, w)
	_, ok := w.(*onlyDifferencesWriter)
	require.True(t, ok)

}

func TestNewOnlyDiffWriterAddMatch(t *testing.T) {

	sb := strings.Builder{}
	s := NewStringWriteCloser(&sb)
	//w := NewOnlyDiffWriter(s).
	//require.NotNil(t, w)
	w, ok := NewOnlyDiffWriter(s).(*onlyDifferencesWriter)
	require.True(t, ok)
	// tcs := []struct{
	// 	desc string
	// 	current_index uint64
	// 	match index uint64
	// 	expected string
	// }{
	// 		{
	// 		"match block 0 on first read",
	// 		0,
	// 		0,
	// 		"",
	// 	},
	// }

	// testcases := []struct {}
	// {
	// 	{},
	// }
	// for _, tc := range tcs{
	// 	t.Run(tc.desc, func(t* testing.T){
	// 		w.blockOffset = tc.current_index
	// 	})
	// }
	w.AddMatch(0)
	w.Flush()
	require.Equal(t, "", sb.String())
	require.Equal(t, uint64(1), w.blockOffset)
	w.AddMatch(1)
	w.Flush()
	require.Equal(t, "", sb.String())
	require.Equal(t, uint64(2), w.blockOffset)
	// imagine block 4 was matched before block 2
	// should write out two removals
	w.AddMatch(4)
	w.Flush()
	require.Equal(t, "- BLOCK_2\n- BLOCK_3\n", sb.String())
	require.Equal(t, uint64(5), w.blockOffset)

	// imagine block 3 was matched after reading block 4
	w.AddMatch(3)
	w.Flush()
	require.Equal(t, "- BLOCK_2\n- BLOCK_3\n= BLOCK_3\n", sb.String())
	require.Equal(t, uint64(5), w.blockOffset)
	// imagine block 4 was matched after reading block 4 & block 4
	w.AddMatch(4)
	w.Flush()
	require.Equal(t, "- BLOCK_2\n- BLOCK_3\n= BLOCK_3\n= BLOCK_4\n", sb.String())
	require.Equal(t, uint64(5), w.blockOffset)

}
