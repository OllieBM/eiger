package hasher

// TODO:
// consider this to make passing a hash.Hash around simpler
// and removing the .sum(notnil)
type Hasher interface {
	Strong([]byte) []byte
	Weak([]byte) uint32
}
