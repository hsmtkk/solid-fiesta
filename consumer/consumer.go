package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hsmtkk/solid-fiesta/env"
	"github.com/hsmtkk/solid-fiesta/waitnats"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	subscribedMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "subscribed_messages",
		Help: "the number of NATS messages subscribed",
	})
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
	exporterPort := env.MandatoryInt("EXPORTER_PORT")

	natsConn := waitnats.WaitNATS(sugar, natsURL)
	defer natsConn.Close()

	handler := newHandler(sugar)
	natsConn.Subscribe(natsSubject, handler.handle)

	exporterAddr := fmt.Sprintf(":%d", exporterPort)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(exporterAddr, nil)
}

func newHandler(sugar *zap.SugaredLogger) *handler {
	return &handler{sugar}
}

type handler struct {
	sugar *zap.SugaredLogger
}

func (hdl *handler) handle(msg *nats.Msg) {
	subscribedMessages.Inc()
	s := string(msg.Data)
	count, err := strconv.Atoi(s)
	if err != nil {
		hdl.sugar.Errorw("failed to parse as int", "message", s, "error", err)
	}
	hdl.sugar.Infow("subscribe", "count", count)
}
