package factory

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/filter"
	"github.com/x0rworld/go-bloomfilter/filter/rotator"
)

type Factory interface {
	NewFilter(ctx context.Context) (filter.Filter, error)
}

type BloomFilterFactory struct {
	cfg config.FactoryConfig
}

// NewFilter returns filters depends on config.FactoryConfig
// if bitmap is not recognized, return BloomFilter with bitmap.Local by default which manipulates in-memory bitset
func (f *BloomFilterFactory) NewFilter(_ context.Context) (filter.Filter, error) {
	var bm bitmap.Bitmap
	switch f.cfg.FilterConfig.BitmapConfig.Type {
	case config.BitmapTypeRedis:
		client := redis.NewClient(&redis.Options{
			Addr:         f.cfg.RedisConfig.Addr,
			ReadTimeout:  f.cfg.RedisConfig.Timeout,
			WriteTimeout: f.cfg.RedisConfig.Timeout,
		})
		bm = bitmap.NewRedis(client, f.cfg.RedisConfig.Key, f.cfg.FilterConfig.M)
	default:
		bm = bitmap.NewLocal(f.cfg.FilterConfig.M)
	}
	return filter.NewBloomFilter(bm, f.cfg.FilterConfig.M, f.cfg.FilterConfig.K), nil
}

type RotatorFactory struct {
	cfg  config.FactoryConfig
	base Factory
}

// NewFilter returns rotator implementing filter that supports doing rotation by goroutine.
func (f *RotatorFactory) NewFilter(ctx context.Context) (filter.Filter, error) {
	return rotator.NewRotator(ctx, f.cfg.RotatorConfig, f.base.NewFilter)
}

// NewFilterFactory does config validation before return Factory based on config.FactoryConfig.
// Return RotatorFactory if rotator is enabled in config.FactoryConfig, otherwise return BloomFilterFactory.
func NewFilterFactory(fCfg config.FactoryConfig) (Factory, error) {
	// validate config
	if err := fCfg.Validate(); err != nil {
		return nil, err
	}

	var factory Factory
	factory = &BloomFilterFactory{cfg: fCfg}

	// wrap BloomFilterFactory if Rotator is enabled
	if fCfg.RotatorConfig.Enable {
		factory = &RotatorFactory{cfg: fCfg, base: factory}
	}
	return factory, nil
}
