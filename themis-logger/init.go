package tlogger

import (
	"github.com/Sirupsen/logrus"
	_ "github.com/infobloxopen/themis/themis-logger/default-logger"
)

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
