package rolling_checksum

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {
	rc := New()

	tcs := []struct {
		in       []byte
		rolling  RollingChecksum
		expected uint32
	}{
		{
			in:       []byte("abc"),
			rolling:  rc,
			expected: 0x249ded0,
		},
		{
			in:       []byte("bcd"),
			rolling:  rc,
			expected: 0x24fde79,
		},
	}

	for _, tc := range tcs {
		out := tc.rolling.Calculate(tc.in)
		exp := tc.expected
		require.Equal(t, exp, out)
	}
}
func TestCalculateNegative(t *testing.T) {
	rc := New()

	tcs := []struct {
		in      []byte
		rolling RollingChecksum
		control []byte
	}{
		{
			in:      []byte("abc"),
			control: []byte("bcd"),
			rolling: rc,
		},
		{
			in:      []byte("abc"),
			control: []byte("abc       "),
			rolling: rc,
		},
	}

	for _, tc := range tcs {
		out := tc.rolling.Calculate(tc.in)

		exp := tc.rolling.Calculate(tc.control)
		require.NotEqual(t, exp, out)
	}
}

func TestRoll(t *testing.T) {
	tcs := []struct {
		in       []byte
		next     byte
		rolling  RollingChecksum
		expected uint32
	}{
		{
			in:       []byte("abc"),
			next:     byte('d'),
			rolling:  New(),
			expected: 0x249ded0,
		},
	}

	for _, tc := range tcs {
		_ = tc.rolling.Calculate(tc.in)
		out := tc.rolling.Roll(tc.in[0], tc.next)
		exp := tc.rolling.Calculate(append(tc.in[1:], tc.next))
		require.Equal(t, exp, out)
	}
}
