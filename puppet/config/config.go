package config

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
)

var cfg *PuppetConfig

type PuppetConfig struct {
	RuntimeMode          int    `json:"runtime-mode"`
	ConsoleAddress       string `json:"console-address"`
	LogLevel             string `json:"log-level"`
	ConsolePort          int    `json:"console-port"`
	Name                 string `json:"name"`
	Directory            string `json:"directory"`
	InstallTimeStamp     int64  `json:"install-time-stamp"`
	SubProcessBufferSize int    `json:"subprocess-buffer-size"`
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

func ApplyGlobalConfig() {
	switch cfg.LogLevel {
	case "trace":
		logrus.SetLevel(LOG_LEVEL_TRACE)
	case "debug":
		logrus.SetLevel(LOG_LEVEL_DEBUG)
	case "info":
		logrus.SetLevel(LOG_LEVEL_INFO)
	case "warn":
		logrus.SetLevel(LOG_LEVEL_WARN)
	case "error":
		logrus.SetLevel(LOG_LEVEL_ERROR)
	case "fatal":
		logrus.SetLevel(LOG_LEVEL_FATAL)
	case "panic":
		logrus.SetLevel(LOG_LEVEL_PANIC)
	default:
		logrus.SetLevel(LOG_LEVEL_INFO)
	}
}
func init() {
	ResetPuppetConfig()
}

func ResetPuppetConfig() {
	cfg = &PuppetConfig{
		RUNTIME_MODE_NORMAL,
		"127.0.0.1",
		"info",
		1133,
		"default-puppet-name",
		"default-puppet-directory-uuid",
		time.Now().Unix(),
		128,
	}
}

func GetInst() *PuppetConfig {
	return cfg
}

const (
	RUNTIME_MODE_DEBUG int = iota
	RUNTIME_MODE_NORMAL
)

const (
	LOG_LEVEL_TRACE = logrus.TraceLevel
	LOG_LEVEL_DEBUG = logrus.DebugLevel
	LOG_LEVEL_INFO  = logrus.InfoLevel
	LOG_LEVEL_WARN  = logrus.WarnLevel
	LOG_LEVEL_ERROR = logrus.ErrorLevel
	LOG_LEVEL_FATAL = logrus.FatalLevel
	LOG_LEVEL_PANIC = logrus.PanicLevel
)
