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
	// setup redis
	client, err := miniredis.Run()
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	// setup factory
	rotateFreq := 5 * time.Second
	ff, err := genFilterFactory(client, rotateFreq)
	if err != nil {
		log.Println(err)
		return
	}

	// setup filter
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bf, err := ff.NewFilter(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	// scenario 1. test dataHello in initialized state
	log.Println("=========== scenario 1 ===========")
	// 1.1. dataHello should not be in filter
	dataHello := "hello"
	exist, _ := bf.Exist(dataHello)
	log.Printf("dataHello: %q, exist: %v\n", dataHello, exist)
	// 1.2. add dataHello to filter
	_ = bf.Add(dataHello)
	log.Printf("add dataHello: %q\n", dataHello)
	// 1.3 data 1 should be in filter
	exist, _ = bf.Exist(dataHello)
	log.Printf("dataHello: %q, exist: %v\n", dataHello, exist)

	// wait for rotation is performed
	log.Printf("sleep %f seconds\n", (rotateFreq + 1*time.Second).Seconds())
	time.Sleep(rotateFreq + 1*time.Second)

	// scenario 2. test dataHello & dataWorld in first rotation
	log.Println("=========== scenario 2 ===========")
	// 2.1. dataWorld should not be in filter
	dataWorld := "world"
	exist, _ = bf.Exist(dataWorld)
	log.Printf("dataWorld: %q, exist: %v\n", dataWorld, exist)
	// 2.2. add dataWorld to filter
	_ = bf.Add(dataWorld)
	log.Printf("add dataWorld: %q\n", dataWorld)
	// 2.3. dataWorld should be in filter
	exist, _ = bf.Exist(dataWorld)
	log.Printf("dataWorld: %q, exist: %v\n", dataWorld, exist)
	// 2.4. dataHello should be kept into filter
	exist, _ = bf.Exist(dataHello)
	log.Printf("dataHello: %q, exist: %v\n", dataHello, exist)

	// wait for rotation is performed
	log.Printf("sleep %f seconds\n", (rotateFreq + 1*time.Second).Seconds())
	time.Sleep(rotateFreq + 1*time.Second)

	// scenario 3. test dataHello and dataWorld in second rotation
	log.Println("=========== scenario 3 ===========")
	// 3.1. dataWorld should be in filter
	exist, _ = bf.Exist(dataWorld)
	log.Printf("dataWorld: %q, exist: %v\n", dataWorld, exist)
	// 3.2. dataHello should not be in filter
	exist, _ = bf.Exist(dataHello)
	log.Printf("dataHello: %q, exist: %v\n", dataHello, exist)
}

func genFilterFactory(client *miniredis.Miniredis, rotateFreq time.Duration) (factory.FilterFactory, error) {
	return factory.NewFilterFactory(
		config.FactoryConfig{
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
			RotatorConfig: config.RotatorConfig{
				Enable: true,
				Freq:   rotateFreq,
			},
		},
	)
}
