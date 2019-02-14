package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Genesis GenesisInfo
	DB      DBInfo
	Log     LogInfo
}

func New() *Config {
	return &Config{}
}

func (p *Config) Init(cfgFile string) error {
	_, err := toml.DecodeFile(cfgFile, p)
	return err
}

type GenesisInfo struct {
	Account string
	Amount  string
}

type DBInfo struct {
	Type string // sqlite3
	Path string // root/.tentermint/data
}

type LogInfo struct {
	Env  string
	Path string
}
