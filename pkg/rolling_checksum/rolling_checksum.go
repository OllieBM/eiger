package rolling_checksum

const MODULO_FACTOR = 65536

type RollingChecksum interface {
}

type adler32 struct {
	m_block_size int
	hash         uint64
	a            uint64
	b            uint64
}

// TODO: move to a new rolling adler32
func New() RollingChecksum {
	return &adler32{}
}

func (r *adler32) Calculate(bytes []byte) uint64 {

	r.hash = 0
	r.a = 0
	r.b = 0
	r.m_block_size = len(bytes)

	for i, c := range bytes {
		r.a += uint64(c)
		r.b += uint64(r.m_block_size-i) * uint64(c)
	}

	r.a = r.a % MODULO_FACTOR
	r.b = r.b % MODULO_FACTOR
	r.hash = r.a + MODULO_FACTOR*r.b

	return r.hash

}

//
func (r *adler32) Roll(out, in byte) uint64 {

	r.a = (r.a - uint64(out) + uint64(in)) % MODULO_FACTOR
	r.b = (r.b - uint64(r.m_block_size)*uint64(out) + r.a) % MODULO_FACTOR
	r.hash = r.a + MODULO_FACTOR*r.b

	return r.hash
}
