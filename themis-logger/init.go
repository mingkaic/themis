package tlogger

import (
	"github.com/Sirupsen/logrus"
	logutil "github.com/infobloxopen/themis/themis-logger/default-logger"
)

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

func Config() interface{} {
	return logutil.Config()
}
