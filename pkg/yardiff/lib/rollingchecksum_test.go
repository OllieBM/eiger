package lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRollingCheckSum(t *testing.T) {

	rcs := RollingCheckSum{}

	hash := rcs.Calculate([]byte("abcdef"))
	hash_2 := rcs.Roll('a', 'h')
	fmt.Println(hash)
	fmt.Println(hash_2)
	require.NotEqual(t, hash, hash_2)
	hash = rcs.Calculate([]byte("bcdefh"))
	require.Equal(t, hash, hash_2)
}
