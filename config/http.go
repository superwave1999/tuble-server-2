package config

import (
	"os"
	"strconv"
)

type HttpConfig struct {
	FrontendHost   string
	TimeoutSeconds uint64
	MaxSize        string
}

func GetHttpConfig() HttpConfig {
	config := HttpConfig{
		FrontendHost:   "*",
		TimeoutSeconds: 4,
		MaxSize:        "3K",
	}
	v := os.Getenv("HTTP_ALLOWED_FRONTEND_HOST")
	if v != "" {
		config.FrontendHost = v
	}
	v = os.Getenv("HTTP_REQUEST_TIMEOUT")
	if v != "" {
		config.TimeoutSeconds, _ = strconv.ParseUint(v, 10, 8)
	}
	v = os.Getenv("HTTP_REQUEST_MAX_SIZE")
	if v != "" {
		config.MaxSize = v
	}
	return config
}
