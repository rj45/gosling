package ir

import (
	"math/bits"
	"strconv"
)

type RegID uint8

func (id RegID) String() string {
	return "r" + strconv.FormatUint(uint64(id), 10)
}

type RegMask uint32

// NewRegMask creates a new register mask.
func NewRegMask(reg ...RegID) RegMask {
	var mask RegMask
	for _, r := range reg {
		mask.SetReg(r)
	}
	return mask
}

// Regs returns the registers in the mask.
func (m RegMask) Regs() []RegID {
	var regs []RegID
	bitmask := uint(m)
	for bitmask != 0 {
		// Find the position of the rightmost set bit
		pos := bits.TrailingZeros(bitmask)
		regs = append(regs, RegID(pos))

		// Unset the rightmost set bit
		bitmask &= bitmask - 1
	}
	return regs
}

// SetReg sets the bit for the given register in the mask.
func (m *RegMask) SetReg(reg RegID) {
	*m |= 1 << reg
}

// ClearReg clears the bit for the given register in the mask.
func (m *RegMask) ClearReg(reg RegID) {
	*m &= ^(1 << reg)
}

// Pop pops the rightmost set bit from the mask and returns it.
func (m *RegMask) Pop() RegID {
	// Find the position of the rightmost set bit
	pos := bits.TrailingZeros(uint(*m))

	// Unset the rightmost set bit
	*m &= *m - 1

	return RegID(pos)
}

// Peek returns the rightmost set bit in the mask.
func (m RegMask) Peek() RegID {
	return RegID(bits.TrailingZeros(uint(m)))
}

// HasReg returns true if the given register is set in the mask.
func (m RegMask) HasReg(reg RegID) bool {
	return m&(1<<reg) != 0
}

// IsEmpty returns true if the mask is empty.
func (m RegMask) IsEmpty() bool {
	return m == 0
}

// Overlaps returns true if the two masks have any registers in common.
func (m RegMask) Overlaps(other RegMask) bool {
	return m&other != 0
}

// Union returns the union of the two masks.
func (m RegMask) Union(other RegMask) RegMask {
	return m | other
}

// Intersect returns the intersection of the two masks.
func (m RegMask) Intersect(other RegMask) RegMask {
	return m & other
}

// String returns the string representation of the mask.
func (m RegMask) String() string {
	var s string
	mask := m
	for !mask.IsEmpty() {
		s += mask.Pop().String()
	}
	return s
}

const (
	R0 RegID = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
)
