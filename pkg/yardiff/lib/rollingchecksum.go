package lib

const MODULO_FACTOR = uint64(1190494759) // selected from http://www.primos.mat.br/2T_en.html
//uint64(1 << 16)

type RollingCheckSum struct {
	m_r          uint64
	m_r1         uint64
	m_r2         uint64
	m_block_size int
}

func (rc *RollingCheckSum) Calculate(data []byte) uint64 {
	rc.m_r = uint64(0)
	rc.m_r1 = uint64(0)
	rc.m_r2 = uint64(0)
	rc.m_block_size = len(data)

	for i, _ := range data {
		rc.m_r1 += uint64(data[i])
		rc.m_r2 += (uint64(rc.m_block_size) - uint64(i)) * uint64(data[i])
	}

	rc.m_r1 = rc.m_r1 % MODULO_FACTOR
	rc.m_r2 = rc.m_r2 % MODULO_FACTOR
	rc.m_r = rc.m_r1 + MODULO_FACTOR*rc.m_r2

	return rc.m_r
}

func (rc *RollingCheckSum) Roll(outgoing byte, incoming byte) uint64 {
	rc.m_r1 = (rc.m_r1 - uint64(outgoing) + uint64(incoming)) % MODULO_FACTOR
	rc.m_r2 = (rc.m_r2 - (uint64(rc.m_block_size) * uint64(outgoing)) + rc.m_r1) % MODULO_FACTOR
	rc.m_r = rc.m_r1 + MODULO_FACTOR*rc.m_r2
	return rc.m_r
}
