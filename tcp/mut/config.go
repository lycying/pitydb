package mut

import (
	"errors"
	"time"
)

type TcpConfig struct {
	TcpKeepAlive bool
	TcpLinger    int
	TcpNoDelay   bool
}
type Config struct {
	TcpConfig

	ReadBufferSize  int
	WriteBufferSize int
	PendingWriteNum int
	BatchWriteNum   int
	Timeout         int64
	AsyncMode       bool

	codec   Codec
	handler Handler
}

func DefaultConfig() *Config {
	cfg := new(Config)

	cfg.TcpKeepAlive = true
	cfg.TcpLinger = 0

	cfg.ReadBufferSize = 1024 * 4
	cfg.WriteBufferSize = 1024 * 4
	cfg.PendingWriteNum = 100
	cfg.BatchWriteNum = 10
	cfg.Timeout = int64(time.Second * 5)
	cfg.AsyncMode = true
	cfg.handler = nil
	cfg.codec = nil
	return cfg
}

func (cfg *Config) validate() error {
	if nil == cfg.handler {
		return errors.New("You sholud supply a callback")
	}

	if nil == cfg.codec {
		return errors.New("You sholud supply a codec")
	}
	return nil
}

func (cfg *Config) SetCodec(c Codec) {
	cfg.codec = c
}
func (cfg *Config) SetCallback(c Handler) {
	cfg.handler = c
}
