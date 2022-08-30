package bitmap

import (
	"github.com/bits-and-blooms/bitset"
)

type InMemory struct {
	bs *bitset.BitSet
	m  uint64
}

func (im *InMemory) CheckBits(locs []uint64) (bool, error) {
	for _, loc := range locs {
		if !im.bs.Test(uint(loc % im.m)) {
			return false, nil
		}
	}
	return true, nil
}

func (im *InMemory) SetBits(locs []uint64) error {
	for _, loc := range locs {
		im.bs.Set(uint(loc % im.m))
	}
	return nil
}

// NewInMemory returns in-memory bitmap which is backed by github.com/bits-and-blooms/bitset.
func NewInMemory(m uint64) *InMemory {
	return &InMemory{
		bs: bitset.New(uint(m)),
		m:  m,
	}
}
