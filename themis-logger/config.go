package tlogger

import (
	"flag"

	"github.com/Sirupsen/logrus"
)

type defaultLogConfig struct {
	LogLevel string
}

var globalConfig LogConfig

func (cfg *defaultLogConfig) lookVar(key string) interface{} {
	if key == "LogLevel" {
		return &cfg.LogLevel
	}
	return nil
}

func (cfg *defaultLogConfig) afterFlagParse() {
	parsedLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatalf("-log-level [debug|info|warning|error|fatal|panic] is supported, got %s", cfg.LogLevel)
	}
	logrus.SetLevel(parsedLevel)
}

func init() {
	defaultCfg := &defaultLogConfig{}
	flag.StringVar(&defaultCfg.LogLevel, "log-level", "info", "[debug|info|warning|error|fatal|panic] set log level")
	if cfg := reconfig(defaultCfg); cfg != nil {
		globalConfig = cfg
	} else {
		globalConfig = defaultCfg
	}
}

// ConfigLog ... configures and returns options pertaining to logs
func ConfigLog() LogConfig {
	globalConfig.afterFlagParse()
	return globalConfig
}
