package factory

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
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

// NewBitmap returns bitmap.Redis. If rotator is enabled, additionally set ttl to redis server via bitmap.RedisOption.
func (rf *RedisBitmapFactory) NewBitmap(ctx context.Context) (bitmap.Bitmap, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         rf.cfg.RedisConfig.Addr,
		ReadTimeout:  rf.cfg.RedisConfig.Timeout,
		WriteTimeout: rf.cfg.RedisConfig.Timeout,
	})
	if rf.cfg.RotatorConfig.Enable {
		// Set ttl with 2 times of freq plus 5 minutes for each bitmap of redis.
		// Rationale:
		// - 2 times of freq: Each bitmap of redis would stay 2 times of freq due to rotation (with current & next).
		// - additional 5 minutes: Basically, it just needs 2 times of freq for rotation.
		//                         However, set gracefully additional 5 minutes here is preventing corner case just in case.
		//                         For example, the bitmap of redis calls SetBits to operate expired bitset deleted by redis server before the rotation is performed.
		return bitmap.NewRedis(ctx, client, rf.cfg.RedisConfig.Key, rf.cfg.FilterConfig.M, bitmap.RedisSetExpireTTL(rf.cfg.RotatorConfig.Freq*2+RedisGracefulExpireTTL))
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
