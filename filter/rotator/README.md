# Rotator (Rotation)

## Supporting bitmaps

- [In-Memory]
- [Redis]

## Rotator Mode

- `default`: perform rotation by `RotatorConfig.Freq`. It's **effective with all bitmaps**.
    - if rotation is performed when `2022-01-02T03:04:05` and `Freq` is `3h`, the next rotation will
      be `2022-01-02T06:04:05`.
- `truncated-time`: perform rotation by truncated time and `RotatorConfig.Freq`. It's **only effective
  with** `bitmap.Redis{}`.
- Refer following [Explanation](#Explanation) for more information.

### Explanation

- Assume we have `config.FactoryConfig` as following:

```go
package main

import (
	"github.com/x0rworld/go-bloomfilter/config"
)

var cfg = config.FactoryConfig{
	FilterConfig: config.FilterConfig{
		BitmapConfig: config.BitmapConfig{Type: config.BitmapTypeRedis},
		M:            100,
		K:            2,
	},
	RedisConfig: config.RedisConfig{
		Addr:    "elasticache.us-west-2.amazonaws.com:6379",
		Timeout: 5 * time.Second,
		Key:     "bloomfilter",
	},
	RotatorConfig: config.RotatorConfig{
		Enable: true,
		Freq:   "3h",
		Mode:   config.RotatorModeTruncatedTime,
	},
}
```

- The current time is 2022-09-06 08:24:31.35128. (timestamp by nano is 1662452671351280000)

Then, the Redis key for bitmap will be generated and manipulated as following:

| Mode             | Redis Key [1]                       | Next Rotation will be at |
|------------------|-------------------------------------|--------------------------|
| `default`        | bloomfilter_1662452671351280000     | 2022-09-06 11:24:31      |
| `truncated-time` | bloomfilter_1662444000000000000 [2] | 2022-09-06 09:00:00      |

- [1]: [bitmap.Redis] manipulates bitmap based on Redis Key.
- [2]: 1662444000000000000 is timestamp by nano, which is at 2022-09-06 06:00:00.

### Note

Ideally, we assume all machines can produce the closed system time, so go-bloomfilter adopts system time for
synchronization with `truncated-time`.
However, it's not robust because we can't guarantee for the consistency of system time and this package doesn't provide
the synchronization for the time yet.

Unless you handle the consistency problem with system time, please consider suitability of your project before you
adopt `truncated-time` way if you'd like to rely on system time for synchronization.

[In-Memory]: ../../bitmap/memory.go

[Redis]: ../../bitmap/redis.go

[bitmap]: ../../bitmap

[bitmap.Redis]: ../../bitmap/redis.go

[config.go]: ../../config/config.go
