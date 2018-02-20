package main

import (
	"encoding/json"
	"flag"
	"math"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/infobloxopen/themis/pdp/ast"
	"github.com/infobloxopen/themis/pdpserver/server"
)

const (
	policyFormatNameYAML = "yaml"
	policyFormatNameJSON = "json"
)

var policyParsers = map[string]ast.Parser{
	policyFormatNameYAML: ast.NewYAMLParser(),
	policyFormatNameJSON: ast.NewJSONParser(),
}

type config struct {
	policy       string
	policyParser ast.Parser
	content      stringSet
	serviceEP    string
	controlEP    string
	tracingEP    string
	healthEP     string
	profilerEP   string
	mem          server.MemLimits
	maxStreams   uint
}

type stringSet []string

func (s *stringSet) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSet) Set(v string) error {
	*s = append(*s, v)
	return nil
}

var conf config

func init() {
	flag.StringVar(&conf.policy, "p", "", "policy file to start with")
	policyFmt := flag.String("pfmt", policyFormatNameYAML, "policy data format \"yaml\" or \"json\"")
	flag.Var(&conf.content, "j", "JSON content files to start with")
	flag.StringVar(&conf.serviceEP, "l", ":5555", "listen for decision requests on this address:port")
	flag.StringVar(&conf.controlEP, "c", ":5554", "listen for policies on this address:port")
	flag.StringVar(&conf.tracingEP, "t", "", "OpenZipkin tracing endpoint")
	flag.StringVar(&conf.healthEP, "health", "", "health check endpoint")
	flag.StringVar(&conf.profilerEP, "pprof", "", "performance profiler endpoint")
	limit := flag.Uint64("mem-limit", 0, "memory limit in megabytes")
	flag.UintVar(&conf.maxStreams, "max-streams", 0, "maximum number of parallel gRPC streams (0 - use gRPC default)")

	flag.Parse()

	p, ok := policyParsers[strings.ToLower(*policyFmt)]
	if !ok {
		log.WithField("format", *policyFmt).Fatal("unknow policy format")
	}
	conf.policyParser = p

	mem, err := server.MakeMemLimits(*limit*1024*1024, 90, 70, 30, 30)
	if err != nil {
		log.WithError(err).Panic("wrong memory limits")
	}
	conf.mem = mem

	if conf.maxStreams > math.MaxUint32 {
		log.WithFields(log.Fields{
			"max-streams": conf.maxStreams,
			"limit":       math.MaxUint32,
		}).Fatal("too big maximum number of parallel gRPC streams")
	}

	block, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.Errorf("error: %s", err)
	}
	log.Infof("Configuration: \n%s", string(block))
}
