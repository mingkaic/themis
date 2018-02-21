package main

import (
	_ "net/http/pprof"
	"runtime"

	"github.com/Sirupsen/logrus"
	log "github.com/infobloxopen/themis/themis-logger"

	"github.com/infobloxopen/themis/pdpserver/server"
)

func main() {
	logger := log.StandardLogger()
	logger.Info("Starting PDP server")

	pdp := server.NewServer(
		server.WithLogger(logger),
		server.WithPolicyParser(conf.policyParser),
		server.WithServiceAt(conf.ServiceEP),
		server.WithControlAt(conf.ControlEP),
		server.WithHealthAt(conf.HealthEP),
		server.WithProfilerAt(conf.ProfilerEP),
		server.WithTracingAt(conf.TracingEP),
		server.WithMemLimits(conf.mem),
		server.WithMaxGRPCStreams(uint32(conf.MaxStreams)),
	)

	err := pdp.LoadPolicies(conf.Policy)
	if err != nil {
		logger.WithFields(
			logrus.Fields{
				"policy": conf.Policy,
				"err":    err,
			},
		).Error("Failed to load policy. Continue with no policy...")
	}

	err = pdp.LoadContent(conf.content)
	if err != nil {
		logger.WithField("err", err).Error("Failed to load content. Continue with no content...")
	}

	runtime.GC()

	err = pdp.Serve()
	if err != nil {
		logger.WithError(err).Error("Failed to run server")
	}
}
