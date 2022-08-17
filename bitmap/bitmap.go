package bitmap

//go:generate mockgen -package mock -destination ../mock/bitmap_mock.go -source=./bitmap.go

type Bitmap interface {
	// CheckBits returns true if all bits on locs have set
	CheckBits(locs []uint64) (bool, error)
	// SetBits sets all bits on locs
	SetBits(locs []uint64) error
}
