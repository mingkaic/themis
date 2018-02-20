package tlogger

import (
	"github.com/Sirupsen/logrus"
	_ "github.com/mingkaic/MockPDP/themis/themis-logger/default-plugin"
)

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
