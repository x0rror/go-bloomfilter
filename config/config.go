// Package config is used to generate filter by factory.
package config

import (
	"errors"
	"fmt"
	"time"
)

const (
	BitmapTypeInMemory BitmapType = "in-memory"
	BitmapTypeRedis    BitmapType = "redis"
)

var (
	defaultFactoryConfig = FactoryConfig{
		FilterConfig: FilterConfig{
			BitmapConfig: BitmapConfig{Type: BitmapTypeInMemory},
			M:            256 * 1024 * 1024, // 256MiB
			K:            3,
		},
	}
	ErrInvalidBitmapType = errors.New("invalid bitmap type")
)

type BitmapType string

func (b BitmapType) Validate() error {
	switch b {
	case BitmapTypeInMemory, BitmapTypeRedis:
		return nil
	}
	return ErrInvalidBitmapType
}

func NewDefaultFactoryConfig() FactoryConfig {
	return defaultFactoryConfig
}

type FilterConfig struct {
	BitmapConfig BitmapConfig
	// M is the number of bit in bloom filter.
	M uint64
	// K is the number of hash function.
	K uint64
}

func (c FilterConfig) Validate() error {
	if err := c.BitmapConfig.Validate(); err != nil {
		return err
	}
	if c.M == 0 {
		return fmt.Errorf("invalid M: %v", c.M)
	}
	if c.K == 0 {
		return fmt.Errorf("invalid K: %v", c.K)
	}
	return nil
}

type BitmapConfig struct {
	Type BitmapType
}

func (b BitmapConfig) Validate() error {
	return b.Type.Validate()
}

type RedisConfig struct {
	Addr    string
	Timeout time.Duration
	Key     string
}

func (c RedisConfig) Validate() error {
	if c.Addr == "" {
		return errors.New("empty addr")
	}
	if c.Key == "" {
		return errors.New("empty key")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout <= 0")
	}
	return nil
}

type RotatorConfig struct {
	Enable bool
	Freq   time.Duration
}

func (c RotatorConfig) Validate() error {
	if c.Freq <= 0 {
		return errors.New("freq <= 0")
	}
	return nil
}

type FactoryConfig struct {
	FilterConfig  FilterConfig
	RedisConfig   RedisConfig
	RotatorConfig RotatorConfig
}

func (c FactoryConfig) Validate() error {
	if err := c.FilterConfig.Validate(); err != nil {
		return err
	}
	if c.FilterConfig.BitmapConfig.Type == BitmapTypeRedis {
		if err := c.RedisConfig.Validate(); err != nil {
			return err
		}
	}
	if c.RotatorConfig.Enable {
		if err := c.RotatorConfig.Validate(); err != nil {
			return err
		}
	}
	return nil
}
