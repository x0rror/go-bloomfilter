package factory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/filter"
	"github.com/x0rworld/go-bloomfilter/filter/rotator"
	"testing"
	"time"
)

func TestNewFilterFactory(t *testing.T) {
	cfg := config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeBitSet,
			},
			M: 100,
			K: 3,
		},
	}
	f, err := NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &BloomFilterFactory{}, f)

	cfg = config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeBitSet,
			},
			M: 100,
			K: 3,
		},
		RotatorConfig: config.RotatorConfig{
			Enable: true,
			Freq:   3 * time.Second,
		},
	}
	f, err = NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &RotatorFactory{}, f)

	// invalid config
	cfg = config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeBitSet,
			},
			M: 100,
			K: 0, // invalid K
		},
		RotatorConfig: config.RotatorConfig{
			Enable: true,
			Freq:   3 * time.Second,
		},
	}
	f, err = NewFilterFactory(cfg)
	assert.Error(t, err)
	assert.Nil(t, f)
}

func TestBloomFilterFactory_NewFilter(t *testing.T) {
	// bitmap: bitset
	cfg := config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeBitSet,
			},
			M: 100,
			K: 3,
		},
	}
	ff, err := NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &BloomFilterFactory{}, ff)

	f, err := ff.NewFilter(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &filter.BloomFilter{}, f)
	bf := f.(*filter.BloomFilter)
	assert.IsType(t, &bitmap.Local{}, bf.BitMap)

	// bitmap: redis
	cfg = config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeRedis,
			},
			M: 100,
			K: 3,
		},
		RedisConfig: config.RedisConfig{
			Addr:    "localhost:6379",
			Timeout: 3 * time.Second,
			Key:     "test",
		},
	}
	ff, err = NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &BloomFilterFactory{}, ff)

	f, err = ff.NewFilter(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &filter.BloomFilter{}, f)
	bf = f.(*filter.BloomFilter)
	assert.IsType(t, &bitmap.Redis{}, bf.BitMap)

	// rotator: enabled
	cfg = config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeBitSet,
			},
			M: 100,
			K: 3,
		},
		RotatorConfig: config.RotatorConfig{
			Enable: true,
			Freq:   3 * time.Second,
		},
	}
	ff, err = NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &RotatorFactory{}, ff)

	f, err = ff.NewFilter(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &rotator.Rotator{}, f)
	r := f.(*rotator.Rotator)
	assert.IsType(t, &filter.BloomFilter{}, r.Current)
	bf = r.Current.(*filter.BloomFilter)
	assert.IsType(t, &bitmap.Local{}, bf.BitMap)
}
