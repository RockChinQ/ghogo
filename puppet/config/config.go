package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

var cfg *PuppetConfig

type PuppetConfig struct {
	RuntimeMode      int    `json:"runtime-mode"`
	ConsoleAddress   string `json:"console-address"`
	ConsolePort      int    `json:"console-port"`
	Name             string `json:"name"`
	Directory        string `json:"directory"`
	InstallTimeStamp int64  `json:"install-time-stamp"`
}

func LoadPuppetConfig(configString string) error {
	return json.Unmarshal([]byte(configString), cfg)
}

func OutputPuppetConfig(path string) error {
	bytes, error := json.MarshalIndent(cfg, "", "\t")
	if error != nil {
		return error
	}
	error = ioutil.WriteFile(path, bytes, 0660)
	if error != nil {
		return error
	}
	return nil
}

func init() {
	ResetPuppetConfig()
}

func ResetPuppetConfig() {
	cfg = &PuppetConfig{
		RUNTIME_MODE_NORMAL,
		"127.0.0.1",
		1133,
		"default-puppet-name",
		"default-puppet-directory-uuid",
		int64(time.Now().Unix()) * 1000,
	}
}

func GetInst() *PuppetConfig {
	return cfg
}

const (
	RUNTIME_MODE_DEBUG int = iota
	RUNTIME_MODE_NORMAL
)
