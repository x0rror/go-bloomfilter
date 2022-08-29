// Package filter manipulates bitmap to check and add the element.
package filter

import "github.com/x0rworld/go-bloomfilter/bitmap"

type Filter interface {
	// GetBitmap returns what kind of bitmap.Bitmap used by Filter
	GetBitmap() bitmap.Bitmap
	// Exist returns whether the data is in bitmap.Bitmap
	Exist(data string) (bool, error)
	// Add adds data into bitmap.Bitmap
	Add(data string) error
}
