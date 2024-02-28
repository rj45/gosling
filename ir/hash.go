package ir

// super fast hash function, modified to take 32 bit integers as input
// http://www.azillionmonkeys.com/qed/hash.html
func hashCode(opcode uint32, inputs ...NodeID) uint32 {
	h := uint32(len(inputs) + 1)

	h += opcode & 0xffff
	tmp := ((opcode >> 16) << 11) ^ h
	h = (h << 16) ^ tmp
	h += h >> 11

	for _, id := range inputs {
		h += uint32(id) & 0xffff
		tmp := ((uint32(id) >> 16) << 11) ^ h
		h = (h << 16) ^ tmp
		h += h >> 11
	}

	h ^= h << 3
	h += h >> 5
	h ^= h << 4
	h += h >> 17
	h ^= h << 25
	h += h >> 6

	return h
}
