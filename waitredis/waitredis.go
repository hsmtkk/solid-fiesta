package waitredis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func WaitRedis(sugar *zap.SugaredLogger, redisHost string, redisPort int) *redis.Client {
	addr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	clt := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	for {
		err := clt.Ping(context.Background()).Err()
		if err == nil {
			sugar.Infow("redis ping success", "host", redisHost, "port", redisPort)
			break
		} else {
			sugar.Infow("redis ping fail", "host", redisHost, "port", redisPort, "error", err)
			time.Sleep(1 * time.Second)
		}
	}
	return clt
}
