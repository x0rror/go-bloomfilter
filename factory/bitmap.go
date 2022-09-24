package factory

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/core"
	"time"
)

var RedisGracefulExpireTTL = 5 * time.Minute

type InMemoryBitmapFactory struct {
	cfg config.FactoryConfig
}

func (imf *InMemoryBitmapFactory) NewBitmap(_ context.Context) (bitmap.Bitmap, error) {
	return bitmap.NewInMemory(imf.cfg.FilterConfig.M), nil
}

type RedisBitmapFactory struct {
	cfg config.FactoryConfig
}

// NewBitmap returns bitmap.Redis.
// If rotator is disabled, return default bitmap.Redis that its key is defined as config and without setting TTL.
//
// In contrast, if rotator is enabled:
//  1. append timestamp of current time to the key and
//  2. additionally set TTL to redis server via bitmap.RedisOption (TTL would be the 2 times of freq plus 5 minutes)
//
// Rationale:
//
//   - 2 times of freq: Each bitmap of redis would stay 2 times of freq due to rotation (being current & next).
//   - additional 5 minutes: Basically, it just needs 2 times of freq for rotation.
//     However, set gracefully additional 5 minutes here is preventing corner case just in case.
//     For example, the bitmap of redis calls SetBits to operate expired bitset deleted by redis server before the rotation is performed.
//
// Besides, it refers to value of context.Context (type is core.BitmapFactoryCtxValue) to generate key of bitmap and set expiry.
//
// For example, key: `go-bloomfilter`, current time is `2022-09-06 08:24:31.35128`; freq is `3h`:
//
// if 1) value.IsRotatorEnabled == true, 2) value.RotatorMode == config.RotatorModeTruncatedTime
//
// 3-1) value.IsNextFilter == false, the key of bitmap would be `go-bloomfilter_1662444000000000000`. (`1662444000000000000` is unix timestamp of `2022-09-06 06:00:00`.)
//
// 3-2) value.IsNextFilter == true, the key of bitmap would be `go-bloomfilter_1662454800000000000`. (`1662454800000000000` is unix timestamp of `2022-09-06 09:00:00`.)
func (rf *RedisBitmapFactory) NewBitmap(ctx context.Context) (bitmap.Bitmap, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         rf.cfg.RedisConfig.Addr,
		ReadTimeout:  rf.cfg.RedisConfig.Timeout,
		WriteTimeout: rf.cfg.RedisConfig.Timeout,
	})
	val, ok := ctx.Value(core.BitmapFactoryCtxKey).(core.BitmapFactoryCtxValue)
	if ok && val.IsRotatorEnabled {
		now := val.Now
		if val.IsNextFilter {
			now = now.Add(rf.cfg.RotatorConfig.Freq)
		}
		if val.RotatorMode == config.RotatorModeTruncatedTime {
			now = now.Truncate(rf.cfg.RotatorConfig.Freq)
		}
		return bitmap.NewRedis(
			ctx,
			client,
			fmt.Sprintf("%s_%d", rf.cfg.RedisConfig.Key, now.UnixNano()),
			rf.cfg.FilterConfig.M,
			bitmap.RedisSetExpireTTL(rf.cfg.RotatorConfig.Freq*2+RedisGracefulExpireTTL),
		)
	} else {
		return bitmap.NewRedis(ctx, client, rf.cfg.RedisConfig.Key, rf.cfg.FilterConfig.M)
	}
}

// NewBitmapFactory does config validation with config.FactoryConfig before returns BitmapFactory depending on cfg.FilterConfig.BitmapConfig.Type.
// If type of bitmap is not recognized, return bitmap.InMemory by default.
func NewBitmapFactory(cfg config.FactoryConfig) (BitmapFactory, error) {
	// validate config
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	switch cfg.FilterConfig.BitmapConfig.Type {
	case config.BitmapTypeRedis:
		return &RedisBitmapFactory{cfg: cfg}, nil
	default:
		return &InMemoryBitmapFactory{cfg: cfg}, nil
	}
}
