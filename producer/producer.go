package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger: %s", err)
	}
	defer logger.Sync()
}
