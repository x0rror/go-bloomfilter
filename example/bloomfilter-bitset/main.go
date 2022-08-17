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
