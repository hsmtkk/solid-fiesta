package env

import (
	"log"
	"os"
	"strconv"
)

func MandatoryString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("you must define %s environment variable", key)
	}
	return val
}

func MandatoryInt(key string) int {
	val := MandatoryString(key)
	n, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("failed to parse %s as int; %s", val, err)
	}
	return n
}
