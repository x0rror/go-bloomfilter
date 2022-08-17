package bitmap

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	client *redis.Client
	key    string
	m      uint64
}

func (r *Redis) CheckBits(locs []uint64) (bool, error) {
	pl := r.client.Pipeline()
	ctx := context.Background()

	var results []*redis.IntCmd
	for _, loc := range locs {
		results = append(results, pl.GetBit(ctx, r.key, int64(loc%r.m)))
	}
	_, err := pl.Exec(ctx)
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
	ctx := context.Background()
	var results []*redis.IntCmd
	for _, loc := range locs {
		results = append(results, pl.SetBit(ctx, r.key, int64(loc%r.m), 1))
	}
	_, err := pl.Exec(ctx)
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

func (r *Redis) SetExpireTTL(d time.Duration) error {
	res := r.client.Expire(context.Background(), r.key, d)
	_, err := res.Result()
	if err != nil {
		return err
	}
	return nil
}

func NewRedis(client *redis.Client, key string, m uint64) *Redis {
	r := &Redis{
		client: client,
		key:    fmt.Sprintf("%s_%d", key, time.Now().UnixNano()),
		m:      m,
	}
	// Set the empty bitmap with the key in Redis.
	// Also, when Rotator applies Redis as bitmap, it will cause failure with setting expiry if the key is not in Redis in prior.
	r.client.SetBit(context.Background(), r.key, 0, 0)
	return r
}
