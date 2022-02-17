package waitnats

import (
	"time"

	"github.com/nats-io/nats.go"
)

func WaitNATS(natsURL string) *nats.Conn {
	var natsConn *nats.Conn
	var err error
	for {
		natsConn, err = nats.Connect(natsURL)
		if err == nil {
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return natsConn
}
