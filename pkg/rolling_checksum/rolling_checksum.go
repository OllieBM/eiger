package rolling_checksum

const MODULO_FACTOR = 65521

type RollingChecksum interface {
	Calculate([]byte) uint32
	Roll(out, in byte) uint32
}

// modified adler32 algorithm
// as used in https://rsync.samba.org/tech_report/node3.html
// hint: no need for the plus 1 because arrays are indexed from 0
type rollingAdler32 struct {
	m_block_size int
	hash         uint32
	a            uint32
	b            uint32
	//buffer 		 []byte//
}

// TODO: move to a new rolling adler32
func New() RollingChecksum {
	return &rollingAdler32{}
}

func (r *rollingAdler32) Calculate(bytes []byte) uint32 {

	r.hash = 0
	r.a = 0
	r.b = 0
	r.m_block_size = len(bytes)

	for i, c := range bytes {
		r.a += uint32(c)
		r.b += uint32(r.m_block_size-i) * uint32(c)
	}

	r.a = r.a % MODULO_FACTOR
	r.b = r.b % MODULO_FACTOR
	r.hash = r.a + MODULO_FACTOR*r.b

	return r.hash

}

//
func (r *rollingAdler32) Roll(out, in byte) uint32 {

	r.a = (r.a - uint32(out) + uint32(in)) % MODULO_FACTOR
	r.b = (r.b - uint32(r.m_block_size)*uint32(out) + r.a) % MODULO_FACTOR
	r.hash = r.a + MODULO_FACTOR*r.b

	return r.hash
}
