package main

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/factory"
	"log"
	"time"
)

func main() {
	client, err := miniredis.Run()
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	ff, err := factory.NewFilterFactory(config.FactoryConfig{
		FilterConfig: config.FilterConfig{
			BitmapConfig: config.BitmapConfig{Type: config.BitmapTypeRedis},
			M:            100,
			K:            2,
		},
		RedisConfig: config.RedisConfig{
			Addr:    client.Addr(),
			Timeout: 5 * time.Second,
			Key:     "filter-redis",
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	bf, err := ff.NewFilter(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	data := "hello world"
	exist, err := bf.Exist(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("data: %v, exist: %v\n", data, exist)
	err = bf.Add(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("add data: %s\n", data)
	exist, err = bf.Exist(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("data: %v, exist: %v\n", data, exist)
	log.Printf("redis key: %s\n", client.Keys())
}
