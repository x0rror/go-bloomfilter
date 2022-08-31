package factory

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
	"regexp"
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

// assertKeyTTL asserts expDur for all matched keys started with prefix `key` (equivalent to `key*`)
func assertKeyTTL(t *testing.T, mr *miniredis.Miniredis, key string, expDur time.Duration) {
	keys := mr.Keys()
	keyPtn, err := regexp.Compile(fmt.Sprintf("%s*", key))
	assert.NoError(t, err)
	for _, k := range keys {
		if keyPtn.MatchString(k) {
			ttl := mr.TTL(k)
			assert.Equal(t, expDur, ttl)
		}
	}
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
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    assert.ErrorAssertionFunc
		wantKeyTTL time.Duration
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
			args:       args{ctx: context.Background()},
			wantErr:    assert.NoError,
			wantKeyTTL: 0,
		},
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
						Key:     "test-RedisBitmapFactory_NewBitmap-rotatorIsEnabled",
					},
					RotatorConfig: config.RotatorConfig{
						Enable: true,
						Freq:   freq,
					},
				},
			},
			args:       args{ctx: context.Background()},
			wantErr:    assert.NoError,
			wantKeyTTL: freq*2 + 5*time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rf := &RedisBitmapFactory{
				cfg: tt.fields.cfg,
			}
			got, err := rf.NewBitmap(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("NewBitmap(%v)", tt.args.ctx)) {
				return
			}
			assert.IsType(t, &bitmap.Redis{}, got)
			assertKeyTTL(t, mr, tt.fields.cfg.RedisConfig.Key, tt.wantKeyTTL)
		})
	}
}
