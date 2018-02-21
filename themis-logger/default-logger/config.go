package defaultlog

import (
	"flag"

	"github.com/Sirupsen/logrus"
)

type defaultLogConfig struct {
	LogLevel string
}

var conf defaultLogConfig

func init() {
	flag.StringVar(&conf.LogLevel, "log-level", "info", "[debug|info|warning|error|fatal|panic] set log level")
}

func Config() interface{} {
	parsedLevel, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		logrus.Fatalf("-log-level [debug|info|warning|error|fatal|panic] is supported, got %s", conf.LogLevel)
	}
	logrus.SetLevel(parsedLevel)

	return &conf
}
