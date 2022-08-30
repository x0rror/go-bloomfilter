// Package filter manipulates bitmap to check and add the element.
package filter

type Filter interface {
	// Exist returns whether the data is in bitmap.Bitmap
	Exist(data string) (bool, error)
	// Add adds data into bitmap.Bitmap
	Add(data string) error
}
