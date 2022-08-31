package factory

import (
	"context"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/filter"
	"github.com/x0rworld/go-bloomfilter/filter/rotator"
)

type BloomFilterFactory struct {
	cfg config.FactoryConfig
}

// NewFilter returns filters depends on config.FactoryConfig.
func (f *BloomFilterFactory) NewFilter(ctx context.Context) (filter.Filter, error) {
	bmf, err := NewBitmapFactory(f.cfg)
	if err != nil {
		return nil, err
	}
	bm, err := bmf.NewBitmap(ctx)
	if err != nil {
		return nil, err
	}
	return filter.NewBloomFilter(bm, f.cfg.FilterConfig.M, f.cfg.FilterConfig.K), nil
}

type RotatorFactory struct {
	cfg  config.FactoryConfig
	base FilterFactory
}

// NewFilter returns rotator implementing filter that supports doing rotation by goroutine.
func (f *RotatorFactory) NewFilter(ctx context.Context) (filter.Filter, error) {
	return rotator.NewRotator(ctx, f.cfg.RotatorConfig, f.base.NewFilter)
}

// NewFilterFactory does config validation with config.FactoryConfig before returns FilterFactory.
// Returns RotatorFactory if rotator is enabled specified within config.FactoryConfig, otherwise return BloomFilterFactory.
func NewFilterFactory(cfg config.FactoryConfig) (FilterFactory, error) {
	// validate config
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var factory FilterFactory
	factory = &BloomFilterFactory{cfg: cfg}

	// wrap BloomFilterFactory if Rotator is enabled
	if cfg.RotatorConfig.Enable {
		factory = &RotatorFactory{cfg: cfg, base: factory}
	}
	return factory, nil
}
