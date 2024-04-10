package config

import (
	"time"

	"github/Thunderbirrd/pow-tcp-server/pkg/utils"
)

type ServerConfig struct {
	Address   string
	Deadline  time.Duration
	KeepAlive time.Duration
	PowConfig *PowConfig
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}

	utils.EnvToStr(&cfg.Address, "SERVER_ADDRESS", "0.0.0.0:80")
	utils.EnvToDuration(&cfg.Deadline, "SERVER_DEADLINE", time.Second*10)
	utils.EnvToDuration(&cfg.KeepAlive, "KEEP_ALIVE", time.Second*10)

	cfg.PowConfig = newPowConfig()

	return cfg
}

type ClientConfig struct {
	ServerAddress string
	MaxRequest    int
	KeepAlive     time.Duration
	PowConfig     *PowConfig
}

func NewClientConfig() *ClientConfig {
	cfg := &ClientConfig{}

	utils.EnvToStr(&cfg.ServerAddress, "SERVER_ADDRESS", "127.0.0.1:80")
	utils.EnvToInt(&cfg.MaxRequest, "MAX_REQUEST", 30)
	utils.EnvToDuration(&cfg.KeepAlive, "KEEP_ALIVE", time.Second*10)

	cfg.PowConfig = newPowConfig()

	return cfg
}

type PowConfig struct {
	Complexity int
}

func newPowConfig() *PowConfig {
	cfg := &PowConfig{}

	utils.EnvToInt(&cfg.Complexity, "POW_COMPLEXITY", 22)

	return cfg
}
