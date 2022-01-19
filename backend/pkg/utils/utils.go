package utils

import (
	"os"
	"strconv"
)

func GetEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func GetPortEnv(key string, defaultVal uint16) uint16 {
	if val := os.Getenv(key); val != "" {
		integer, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			panic("Error to get Postgres port environment")
		}
		return uint16(integer)
	}
	return defaultVal
}
