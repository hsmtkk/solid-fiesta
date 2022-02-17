package main

import (
	"log"
	"strconv"

	"github.com/hsmtkk/solid-fiesta/env"
	"github.com/nats-io/nats.go"
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

	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		sugar.Fatalw("failed to connect NATS", "URL", natsURL, "error", err)
	}
	defer natsConn.Close()

	handler := newHandler(sugar)

	natsConn.Subscribe(natsSubject, handler.handle)

	select {}
}

func newHandler(sugar *zap.SugaredLogger) *handler {
	return &handler{sugar}
}

type handler struct {
	sugar *zap.SugaredLogger
}

func (hdl *handler) handle(msg *nats.Msg) {
	s := string(msg.Data)
	count, err := strconv.Atoi(s)
	if err != nil {
		hdl.sugar.Errorw("failed to parse as int", "message", s, "error", err)
	}
	hdl.sugar.Infow("subscribe", "count", count)
}
