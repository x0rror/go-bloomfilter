package factory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
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
	bloomFilterFactory, err := NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &BloomFilterFactory{}, bloomFilterFactory)

	filter, err := bloomFilterFactory.NewFilter(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &bitmap.Local{}, filter.GetBitmap())

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
	bloomFilterFactory, err = NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &BloomFilterFactory{}, bloomFilterFactory)

	filter, err = bloomFilterFactory.NewFilter(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &bitmap.Redis{}, filter.GetBitmap())

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
	bloomFilterFactory, err = NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &RotatorFactory{}, bloomFilterFactory)

	filter, err = bloomFilterFactory.NewFilter(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &bitmap.Local{}, filter.GetBitmap())
}
