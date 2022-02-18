package waitnats

import (
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func WaitNATS(sugar *zap.SugaredLogger, natsURL string) *nats.Conn {
	var natsConn *nats.Conn
	var err error
	for {
		natsConn, err = nats.Connect(natsURL)
		if err == nil {
			sugar.Infow("connected NATS", "URL", natsURL)
			break
		} else {
			sugar.Infow("waiting NATS", "URL", natsURL, "error", err)
			time.Sleep(1 * time.Second)
		}
	}
	return natsConn
}
