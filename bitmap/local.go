package bitmap

import (
	"github.com/bits-and-blooms/bitset"
)

type Local struct {
	bs *bitset.BitSet
	m  uint64
}

func (l *Local) CheckBits(locs []uint64) (bool, error) {
	for _, loc := range locs {
		if !l.bs.Test(uint(loc % l.m)) {
			return false, nil
		}
	}
	return true, nil
}

func (l *Local) SetBits(locs []uint64) error {
	for _, loc := range locs {
		l.bs.Set(uint(loc % l.m))
	}
	return nil
}

// NewLocal returns in-memory bitmap which is backed by github.com/bits-and-blooms/bitset.
func NewLocal(m uint64) *Local {
	return &Local{
		bs: bitset.New(uint(m)),
		m:  m,
	}
}
