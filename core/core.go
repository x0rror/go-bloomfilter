package core

import (
	"github.com/x0rworld/go-bloomfilter/config"
	"time"
)

type CtxKey string

var (
	BitmapFactoryCtxKey = CtxKey("bitmap-factory-ctx-key")
)

type BitmapFactoryCtxValue struct {
	IsRotatorEnabled bool
	IsNextFilter     bool
	RotatorMode      config.RotatorMode
	Now              time.Time
}
