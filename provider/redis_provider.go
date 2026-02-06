package provider

import (
	"github.com/SexyBobRiK/gostrap/config"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type RedisProvider struct{}

func (RedisProvider) ProviderInit(entities []config.RedisEntity) (map[int]*redis.Client, error) {
	var (
		dbMap = make(map[int]*redis.Client)
		mu    sync.Mutex
		eg    errgroup.Group
	)
	for _, entity := range entities {
		if !entity.Enable {
			continue
		}
		eg.Go(func() error {
			rdb := redis.NewClient(&redis.Options{
				Addr:     entity.Addr,
				Password: entity.Password,
				DB:       entity.DB,
				Protocol: entity.Protocol,
			})
			mu.Lock()
			dbMap[entity.DB] = rdb
			mu.Unlock()
			log.Printf("[Gostrap] Redis %s connected", entity.DB)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return dbMap, nil
}
