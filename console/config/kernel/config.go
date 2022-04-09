//Manage config of console and addons
package kernel

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

var cfg *KernelConfig

//Config of console main runtime
type KernelConfig struct {
	RuntimeMode   int       `json:"runtime-mode"`
	LogLevel      log.Level `json:"log-level"`
	SocketAddress string    `json:"socket-address"`
	SocketPort    int       `json:"socket-port"`
	PuppetTimeOut int       `json:"puppet-time-out"`
}

//Initialize cfg from a json
func LoadKernelConfig(configString string) error {
	return json.Unmarshal([]byte(configString), cfg)
}

func OutputKernelConfig(path string) error {
	bytes, error := json.MarshalIndent(cfg, "", "\t")
	if error != nil {
		return error
	}
	//Can only be R/W by owner and other users in same group
	error = ioutil.WriteFile(path, bytes, 0660)
	if error != nil {
		return error
	}
	return nil
}

func init() {
	ResetKernelConfig()
}

//Reset cfg to default values
func ResetKernelConfig() {
	cfg = &KernelConfig{
		RUNTIME_MODE_NORMAL,
		LOG_LEVEL_INFO,
		"",
		1133,
		10000,
	}
}

func ApplyGlobalConfig() {
	log.SetLevel(cfg.LogLevel)
}

func GetInst() *KernelConfig {
	return cfg
}

const (
	RUNTIME_MODE_DEBUG int = iota
	RUNTIME_MODE_NORMAL
)

const (
	LOG_LEVEL_TRACE = log.TraceLevel
	LOG_LEVEL_DEBUG = log.DebugLevel
	LOG_LEVEL_INFO  = log.InfoLevel
	LOG_LEVEL_WARN  = log.WarnLevel
	LOG_LEVEL_ERROR = log.ErrorLevel
	LOG_LEVEL_FATAL = log.FatalLevel
	LOG_LEVEL_PANIC = log.PanicLevel
)
