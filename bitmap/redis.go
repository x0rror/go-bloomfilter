package bitmap

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisOption func(*Redis) error

type Redis struct {
	ctx    context.Context
	client *redis.Client
	key    string
	m      uint64
}

func (r *Redis) CheckBits(locs []uint64) (bool, error) {
	pl := r.client.Pipeline()

	var results []*redis.IntCmd
	for _, loc := range locs {
		results = append(results, pl.GetBit(r.ctx, r.key, int64(loc%r.m)))
	}
	_, err := pl.Exec(r.ctx)
	if err != nil {
		return false, err
	}
	for _, v := range results {
		res, err := v.Result()
		if err != nil {
			return false, err
		}
		if res == 0 {
			return false, nil
		}
	}
	return true, nil
}

func (r *Redis) SetBits(locs []uint64) error {
	pl := r.client.Pipeline()
	var results []*redis.IntCmd
	for _, loc := range locs {
		results = append(results, pl.SetBit(r.ctx, r.key, int64(loc%r.m), 1))
	}
	_, err := pl.Exec(r.ctx)
	if err != nil {
		return err
	}
	for _, v := range results {
		_, err := v.Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// RedisSetExpireTTL sets expiry TTL with d.
func RedisSetExpireTTL(d time.Duration) RedisOption {
	return func(r *Redis) error {
		res := r.client.Expire(r.ctx, r.key, d)
		_, err := res.Result()
		if err != nil {
			return err
		}
		return nil
	}
}

func (r *Redis) setEmptyBitmap() error {
	res, err := r.client.Keys(r.ctx, r.key).Result()
	if err != nil {
		return err
	}
	if len(res) == 0 {
		r.client.SetBit(r.ctx, r.key, 0, 0)
	}
	return nil
}

// NewRedis returns bitmap that is store into redis and manipulated via github.com/go-redis/redis.
func NewRedis(ctx context.Context, client *redis.Client, key string, m uint64, opts ...RedisOption) (*Redis, error) {
	r := &Redis{
		ctx:    ctx,
		client: client,
		key:    key,
		m:      m,
	}

	// Set the empty bitmap with the r.key to avoid subsequent RedisOption might not be effective such as RedisSetExpireTTL.
	err := r.setEmptyBitmap()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err := opt(r)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
