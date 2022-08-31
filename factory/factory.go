// Package factory is used to generate instances such as filter, bitmap.
package factory

import (
	"context"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/filter"
)

type BitmapFactory interface {
	// NewBitmap generates bitmap.
	NewBitmap(ctx context.Context) (bitmap.Bitmap, error)
}

type FilterFactory interface {
	// NewFilter generates filter.
	NewFilter(ctx context.Context) (filter.Filter, error)
}
