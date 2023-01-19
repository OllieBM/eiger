package rolling_checksum

// see https://www.rfc-editor.org/rfc/rfc1950
// and the adler32 algorithm to better understand these values
// const MODULO_FACTOR = 65521 // 65521 is the largest prime that can be stored in a 16-bit unsigned integer
// const MODULO_FACTOR = 2147483647
// const NMAX = 5552

//
// TRY THIS : https://github.com/bakergo/rollsum/blob/master/rollsum.go

// modified adler32 algorithm
// as used in https://rsync.samba.org/tech_report/node3.html
// hint: no need for the plus 1 because arrays are indexed from 0
type improvedAdler struct {
	m_block_size int
	//hash         uint32
	a      uint16
	b      uint16
	buffer []byte //
}
