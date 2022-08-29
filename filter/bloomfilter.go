package filter

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/x0rworld/go-bloomfilter/bitmap"
)

// locationFunc returns hash locations based on data and k.
type locationFunc func(data []byte, k uint) []uint64

type BloomFilter struct {
	bitmap bitmap.Bitmap
	// m is the number of bit in bloom filter.
	m uint64
	// k is the number of hash function.
	k        uint64
	location locationFunc
}

func (b *BloomFilter) GetBitmap() bitmap.Bitmap {
	return b.bitmap
}

func (b *BloomFilter) Exist(data string) (bool, error) {
	locs := b.location([]byte(data), uint(b.k))
	exist, err := b.bitmap.CheckBits(locs)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (b *BloomFilter) Add(data string) error {
	locs := b.location([]byte(data), uint(b.k))
	err := b.bitmap.SetBits(locs)
	if err != nil {
		return err
	}
	return nil
}

func NewBloomFilter(bitmap bitmap.Bitmap, m, k uint64) *BloomFilter {
	return &BloomFilter{
		bitmap:   bitmap,
		m:        m,
		k:        k,
		location: bloom.Locations,
	}
}
