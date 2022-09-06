package factory

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/core"
	"testing"
	"time"
)

func TestNewBitmapFactory(t *testing.T) {
	// invalid config: redis addr is empty
	cfg := config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeRedis,
			},
			M: 100,
			K: 3,
		},
		RedisConfig: config.RedisConfig{
			Addr:    "",
			Timeout: 1,
			Key:     "test-TestNewBitmapFactory-emptyRedisAddr",
		},
	}
	bmf, err := NewBitmapFactory(cfg)
	assert.Error(t, err)
	assert.Nil(t, bmf)

	// bitmap: in-memory
	cfg = config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeInMemory,
			},
			M: 100,
			K: 3,
		},
	}
	bmf, err = NewBitmapFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &InMemoryBitmapFactory{}, bmf)

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
			Timeout: time.Second,
			Key:     "test-TestNewBitmapFactory-redis",
		},
	}
	bmf, err = NewBitmapFactory(cfg)
	assert.NoError(t, err)
	assert.IsType(t, &RedisBitmapFactory{}, bmf)
}

func TestInMemoryBitmapFactory_NewBitmap(t *testing.T) {
	cfg := config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{
				Type: config.BitmapTypeInMemory,
			},
			M: 100,
			K: 3,
		},
	}
	imbf := &InMemoryBitmapFactory{cfg: cfg}
	imb, err := imbf.NewBitmap(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &bitmap.InMemory{}, imb)
}

// assertKeyTTL asserts expDur for all matched keys started with `key`
func assertKeyTTL(t *testing.T, mr *miniredis.Miniredis, key string, expDur time.Duration) {
	matched := false
	keys := mr.Keys()
	for _, k := range keys {
		if k == key {
			ttl := mr.TTL(k)
			assert.Equal(t, expDur, ttl)
			matched = true
		}
	}
	if !matched {
		assert.Fail(t, "no matched key to be checked with TTL.")
	}
}

var fakeTimeFunc = func() time.Time {
	t, _ := time.Parse(time.RFC3339, "2020-01-02T03:04:05Z00:00")
	return t
}

func TestRedisBitmapFactory_NewBitmap(t *testing.T) {
	mr := miniredis.RunT(t)
	defer mr.Close()

	freq := 10 * time.Second

	type fields struct {
		cfg config.FactoryConfig
	}
	type args struct {
		ctx context.Context
	}
	type expect struct {
		redisKey    string
		redisKeyTTL time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		// expect.redisKey & expect.redisKeyTTL will be asserted based on the result of miniredis
		expect struct {
			redisKey    string
			redisKeyTTL time.Duration
		}
	}{
		{
			name: "rotator is disabled",
			fields: fields{
				cfg: config.FactoryConfig{
					FilterConfig: config.FilterConfig{
						BitmapConfig: config.BitmapConfig{
							Type: config.BitmapTypeRedis,
						},
						M: 100,
						K: 3,
					},
					RedisConfig: config.RedisConfig{
						Addr:    mr.Addr(),
						Timeout: time.Second,
						Key:     "test-RedisBitmapFactory_NewBitmap-rotatorIsDisabled",
					},
					RotatorConfig: config.RotatorConfig{
						Enable: false,
						Freq:   freq,
					},
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: assert.NoError,
			expect: expect{
				redisKey:    "test-RedisBitmapFactory_NewBitmap-rotatorIsDisabled",
				redisKeyTTL: 0,
			},
		},
		{
			name: "rotator is enabled",
			fields: fields{
				cfg: config.FactoryConfig{
					FilterConfig: config.FilterConfig{
						BitmapConfig: config.BitmapConfig{
							Type: config.BitmapTypeRedis,
						},
						M: 100,
						K: 3,
					},
					RedisConfig: config.RedisConfig{
						Addr:    mr.Addr(),
						Timeout: time.Second,
						Key:     "test-RedisBitmapFactory_NewBitmap-rotatorIsEnabled",
					},
					RotatorConfig: config.RotatorConfig{
						Enable: true,
						Freq:   freq,
					},
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: assert.NoError,
			expect: expect{
				redisKey:    fmt.Sprintf("%s_%d", "test-RedisBitmapFactory_NewBitmap-rotatorIsEnabled", fakeTimeFunc().UnixNano()),
				redisKeyTTL: freq*2 + 5*time.Minute,
			},
		},
		{
			name: "rotator is enabled: type = truncated-time",
			fields: fields{
				cfg: config.FactoryConfig{
					FilterConfig: config.FilterConfig{
						BitmapConfig: config.BitmapConfig{
							Type: config.BitmapTypeRedis,
						},
						M: 100,
						K: 3,
					},
					RedisConfig: config.RedisConfig{
						Addr:    mr.Addr(),
						Timeout: time.Second,
						Key:     "test-RedisBitmapFactory_NewBitmap-rotatorIsEnabled-typeIsTruncatedTime",
					},
					RotatorConfig: config.RotatorConfig{
						Enable: true,
						Freq:   freq,
						Mode:   config.RotatorModeTruncatedTime,
					},
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: assert.NoError,
			expect: expect{
				redisKey:    fmt.Sprintf("%s_%d", "test-RedisBitmapFactory_NewBitmap-rotatorIsEnabled-typeIsTruncatedTime", fakeTimeFunc().Truncate(freq).UnixNano()),
				redisKeyTTL: freq*2 + 5*time.Minute,
			},
		},
		{
			name: "rotator is enabled: type = truncated-time; validate next bf",
			fields: fields{
				cfg: config.FactoryConfig{
					FilterConfig: config.FilterConfig{
						BitmapConfig: config.BitmapConfig{
							Type: config.BitmapTypeRedis,
						},
						M: 100,
						K: 3,
					},
					RedisConfig: config.RedisConfig{
						Addr:    mr.Addr(),
						Timeout: time.Second,
						Key:     "test-RedisBitmapFactory_NewBitmap-rotatorIsEnabled-typeIsTruncatedTime-validateNextBf",
					},
					RotatorConfig: config.RotatorConfig{
						Enable: true,
						Freq:   freq,
						Mode:   config.RotatorModeTruncatedTime,
					},
				},
			},
			args: args{ctx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, core.ContextKeyFactoryIsNextBm, true)
				return ctx
			}()},
			wantErr: assert.NoError,
			expect: expect{
				redisKey:    fmt.Sprintf("%s_%d", "test-RedisBitmapFactory_NewBitmap-rotatorIsEnabled-typeIsTruncatedTime-validateNextBf", fakeTimeFunc().Add(freq).Truncate(freq).UnixNano()),
				redisKeyTTL: freq*2 + 5*time.Minute,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rf := &RedisBitmapFactory{
				cfg: tt.fields.cfg,
				now: fakeTimeFunc(),
			}
			got, err := rf.NewBitmap(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("NewBitmap(%v)", tt.args.ctx)) {
				return
			}
			assert.IsType(t, &bitmap.Redis{}, got)
			assertKeyTTL(t, mr, tt.expect.redisKey, tt.expect.redisKeyTTL)
		})
	}
}
