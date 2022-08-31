package factory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/filter"
	"github.com/x0rworld/go-bloomfilter/filter/rotator"
	"testing"
	"time"
)

func TestNewFilterFactory(t *testing.T) {
	// valid: bloomfilter
	cfg := config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeInMemory,
			},
			M: 100,
			K: 3,
		},
	}
	f, err := NewFilterFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &BloomFilterFactory{}, f)

	// valid: rotator filter
	cfg = config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeInMemory,
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
	rf := f.(*RotatorFactory)
	assert.IsType(t, &BloomFilterFactory{}, rf.base)

	// invalid config
	cfg = config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeInMemory,
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
	// error
	ff := &BloomFilterFactory{
		cfg: config.FactoryConfig{
			FilterConfig: config.FilterConfig{
				BitmapConfig: config.BitmapConfig{
					Type: config.BitmapTypeInMemory,
				},
				M: 0, // invalid
				K: 3,
			},
		},
	}
	f, err := ff.NewFilter(context.Background())
	assert.Error(t, err)
	assert.Nil(t, f)

	// without error
	ff = &BloomFilterFactory{
		cfg: config.FactoryConfig{
			FilterConfig: config.FilterConfig{
				BitmapConfig: config.BitmapConfig{
					Type: config.BitmapTypeInMemory,
				},
				M: 100,
				K: 3,
			},
		},
	}
	f, err = ff.NewFilter(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &filter.BloomFilter{}, f)
}

func TestRotatorFactory_NewFilter(t *testing.T) {
	cfg := config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeInMemory,
			},
			M: 100,
			K: 3,
		},
		RotatorConfig: config.RotatorConfig{
			Enable: true,
			Freq:   3 * time.Second,
		},
	}
	base, err := NewFilterFactory(cfg)
	assert.NoError(t, err)
	rf := &RotatorFactory{
		cfg:  cfg,
		base: base,
	}
	f, err := rf.NewFilter(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &rotator.Rotator{}, f)
}
