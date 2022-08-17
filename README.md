# go-bloomfilter

![build workflow](https://github.com/x0rworld/go-bloomfilter/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/x0rworld/go-bloomfilter.svg)](https://pkg.go.dev/github.com/x0rworld/go-bloomfilter)

[go-bloomfilter] is a bloomfilter implemented by Golang which supports in-memory and Redis. Also, it supports rotation
by period.

## Resources

- [Documentation]
- [Examples]

## Features

- [Supporting bitmap]
    - in-memory: backed by [bits-and-blooms/bloom]
    - redis: backed by [go-redis/redis]
- [Rotation]

## Installation

```shell
go get github.com/x0rworld/go-bloomfilter
```

## Quickstart

```go
package main

import (
	"context"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/factory"
	"log"
)

func main() {
	m, k := bloom.EstimateParameters(100, 0.01)
	// configure factor config
	cfg := config.NewDefaultFactoryConfig()
	// modify config
	cfg.FilterConfig.M = uint64(m)
	cfg.FilterConfig.K = uint64(k)
	// create factory by config
	ff, err := factory.NewFilterFactory(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	// create filter by factory
	f, err := ff.NewFilter(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	// manipulate filter: Exist & Add
	data := "hello world"
	exist, err := f.Exist(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("data: %v, exist: %v\n", data, exist)
	err = f.Add(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("add data: %s\n", data)
	exist, err = f.Exist(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("data: %v, exist: %v\n", data, exist)
}
```

[go-bloomfilter]: https://github.com/x0rworld/go-bloomfilter

[bits-and-blooms/bloom]: https://github.com/bits-and-blooms/bloom

[go-redis/redis]: https://github.com/go-redis/redis

[Examples]: ./example

[Documentation]: https://pkg.go.dev/github.com/x0rworld/go-bloomfilter

[Supporting bitmap]: ./bitmap

[Rotation]: ./filter/rotator