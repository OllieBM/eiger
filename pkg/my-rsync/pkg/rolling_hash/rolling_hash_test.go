package rolling_hash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {

	a := "abc"
	b := "bcd"

	hashA := hash(a, 5)
	hashB := hash(b, 5)

	fmt.Printf("old: %d vs new :%d\n", hashA, hashB)
	require.NotEqual(t, hashA, hashB)
}

func TestRollingHash(t *testing.T) {

	in := "abcdef"

	a := "abc"
	b := "bcd"
	c := "cde"
	d := "def"

	expected := []uint64{hash(a, prime),
		hash(b, prime),
		hash(c, prime),
		hash(d, prime),
	}

	unique := make(map[uint64]struct{})
	for _, h := range expected {
		_, ok := unique[h]
		require.False(t, ok)
		unique[h] = struct{}{}
	}
	hashes := rolling_hash(in, prime)
	fmt.Println(expected, hashes)
	require.Equal(t, expected, hashes)
}
