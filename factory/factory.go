// Package factory is used to generate filter.
package factory

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/filter"
	"github.com/x0rworld/go-bloomfilter/filter/rotator"
	"time"
)

type Factory interface {
	// NewFilter generates filter.
	NewFilter(ctx context.Context) (filter.Filter, error)
}

type BloomFilterFactory struct {
	cfg config.FactoryConfig
}

// NewFilter returns filters depends on config.FactoryConfig.
// If bitmap is not recognized, return BloomFilter with bitmap.Local by default which manipulates in-memory bitset.
func (f *BloomFilterFactory) NewFilter(ctx context.Context) (filter.Filter, error) {
	var bm bitmap.Bitmap
	var err error
	switch f.cfg.FilterConfig.BitmapConfig.Type {
	case config.BitmapTypeRedis:
		client := redis.NewClient(&redis.Options{
			Addr:         f.cfg.RedisConfig.Addr,
			ReadTimeout:  f.cfg.RedisConfig.Timeout,
			WriteTimeout: f.cfg.RedisConfig.Timeout,
		})
		if f.cfg.RotatorConfig.Enable {
			// Set ttl with 2 times of freq plus 5 minutes for each bitmap of redis.
			// Rationale:
			// - 2 times of freq: Each bitmap of redis would stay 2 times of freq due to rotation (with current & next).
			// - additional 5 minutes: Basically, it just needs 2 times of freq for rotation.
			//                         However, set gracefully additional 5 minutes here is preventing corner case just in case.
			//                         For example, the bitmap of redis calls SetBits to operate expired bitset deleted by redis server before the rotation is performed.
			bm, err = bitmap.NewRedis(ctx, client, f.cfg.RedisConfig.Key, f.cfg.FilterConfig.M, bitmap.RedisSetExpireTTL(f.cfg.RotatorConfig.Freq*2+5*time.Minute))
			if err != nil {
				return nil, err
			}
		}
		bm, err = bitmap.NewRedis(ctx, client, f.cfg.RedisConfig.Key, f.cfg.FilterConfig.M)
		if err != nil {
			return nil, err
		}
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

// NewFilterFactory does config validation with config.FactoryConfig before return Factory.
// Returns RotatorFactory if rotator is enabled specified within config.FactoryConfig, otherwise return BloomFilterFactory.
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
