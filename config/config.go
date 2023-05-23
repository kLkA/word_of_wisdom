package config

import (
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Client    ClientSettigns    `toml:"client"`
	Server    ServerSettings    `toml:"server"`
	Handshake HandshakeSettings `toml:"handshake"`
}
type ClientSettigns struct {
	Timeout     int `toml:"timeout"`
	WorkerCount int `toml:"worker_count"  default:"4"`
}
type ServerSettings struct {
	ListenAddr string `toml:"listen_addr"`
}
type HandshakeSettings struct {
	Difficulty   int `toml:"difficulty"`
	PowTokenSize int `toml:"pow_token_size"`
}

func GetConfigFromFile(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err = toml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
