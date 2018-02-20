package tlogger

import (
	"github.com/Sirupsen/logrus"
)

//// plugin for configuration

type LogConfig interface {
	lookVar(string) interface{}
	afterFlagParse()
}

// return replacement to globalConfig
func reconfig(cfg LogConfig) LogConfig {
	return nil
}

//// plugin for hook

type LoggerHook interface {
	Start()
	Stop()
}

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

func StandardHook() LoggerHook {
	return nil
}
