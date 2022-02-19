package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hsmtkk/solid-fiesta/env"
	"github.com/hsmtkk/solid-fiesta/waitredis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	subscribedMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "subscribed_messages",
		Help: "the number of messages subscribed",
	})
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger: %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	redisHost := env.MandatoryString("REDIS_HOST")
	redisPort := env.MandatoryInt("REDIS_PORT")
	redisChannel := env.MandatoryString("REDIS_CHANNEL")
	exporterPort := env.MandatoryInt("EXPORTER_PORT")

	go func() {
		exporterAddr := fmt.Sprintf(":%d", exporterPort)
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(exporterAddr, nil)
	}()

	redisConn := waitredis.WaitRedis(sugar, redisHost, redisPort)
	defer redisConn.Close()

	pubsub := redisConn.Subscribe(context.Background(), redisChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		subscribedMessages.Inc()
		count, err := strconv.Atoi(msg.String())
		if err != nil {
			sugar.Errorw("failed to parse as int", "message", msg, "error", err)
			continue
		}
		sugar.Infow("subscribe", "count", count)
	}
}
