package filter

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/x0rworld/go-bloomfilter/mock"
	"testing"
)

var (
	dataHello = "hello"
	dataNone  = ""

	locationHello = []uint64{1}
	locationNone  = []uint64{0}

	errInternal = errors.New("internal error")
)

func stubLocation(data []byte, _ uint) []uint64 {
	if string(data) == "hello" {
		return locationHello
	}
	return locationNone
}

func TestBloomFilter_Exist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bMap := mock.NewMockBitmap(ctrl)
	bf := NewBloomFilter(bMap, 100, 3)
	bf.location = stubLocation

	var exist bool
	var err error

	// return err
	bMap.EXPECT().CheckBits(locationNone).Return(false, errInternal)
	exist, err = bf.Exist(dataNone)
	assert.Error(t, err)
	assert.Equal(t, false, exist)

	// data is not in bloomfilter
	bMap.EXPECT().CheckBits(locationHello).Return(false, nil)
	exist, err = bf.Exist(dataHello)
	assert.NoError(t, err)
	assert.Equal(t, false, exist)

	// data is in bloomfilter
	bMap.EXPECT().CheckBits(locationHello).Return(true, nil)
	exist, err = bf.Exist(dataHello)
	assert.NoError(t, err)
	assert.Equal(t, true, exist)
}

func TestBloomFilter_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bMap := mock.NewMockBitmap(ctrl)
	bf := NewBloomFilter(bMap, 100, 3)
	bf.location = stubLocation

	var err error

	// return err
	bMap.EXPECT().SetBits(locationNone).Return(errInternal)
	err = bf.Add(dataNone)
	assert.Error(t, err)

	// data is added to bloomfilter without err
	bMap.EXPECT().SetBits(locationHello).Return(nil)
	err = bf.Add(dataHello)
	assert.NoError(t, err)
}
