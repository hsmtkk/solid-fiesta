package main

import (
	"log"
	"strconv"
	"time"

	"github.com/hsmtkk/solid-fiesta/env"
	"github.com/hsmtkk/solid-fiesta/waitnats"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger: %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	natsURL := env.MandatoryString("NATS_URL")
	natsSubject := env.MandatoryString("NATS_SUBJECT")
	intervalSeconds := env.MandatoryInt("INTERVAL_SECONDS")

	natsConn := waitnats.WaitNATS(natsURL)
	defer natsConn.Close()

	count := 0
	for {
		msg := strconv.Itoa(count)
		natsConn.Publish(natsSubject, []byte(msg))
		sugar.Infow("publish", "count", count)
		count += 1
		time.Sleep(time.Duration(intervalSeconds) * time.Second)
	}
}
