package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hsmtkk/solid-fiesta/env"
	"github.com/hsmtkk/solid-fiesta/waitredis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	publishedMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "published_messages",
		Help: "the number of messages published",
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
	intervalSeconds := env.MandatoryInt("INTERVAL_SECONDS")
	exporterPort := env.MandatoryInt("EXPORTER_PORT")

	redisConn := waitredis.WaitRedis(sugar, redisHost, redisPort)
	defer redisConn.Close()

	go func() {
		exporterAddr := fmt.Sprintf(":%d", exporterPort)
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(exporterAddr, nil)
	}()

	count := 0
	for {
		msg := strconv.Itoa(count)
		redisConn.Publish(context.Background(), redisChannel, []byte(msg))
		sugar.Infow("publish", "count", count)
		count += 1
		publishedMessages.Inc()
		time.Sleep(time.Duration(intervalSeconds) * time.Second)
	}
}
