package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"time"
)

type CommandLineArgs struct {
	port       int
	host       string
	verbose    bool
	content    string
	headers    []string
	deadTime   string
	minWait    string
	maxWait    string
	randomWait bool
	minSize    string
	maxSize    string
	randomSize bool
}

type ConfigurationReader interface {
	Read() Configuration
}

type ArgsConfigurationReader struct {
	args        CommandLineArgs
	defaultTime time.Duration
}

func (instance *ArgsConfigurationReader) Read() Configuration {
	config := Configuration{}
	config.port = instance.args.port
	config.host = instance.args.host
	config.verbose = instance.args.verbose
	config.content = instance.args.content
	config.headers = instance.args.headers
	config.deadTime = parseTime(instance.args.deadTime)
	config.minWait = parseTime(instance.args.minWait)
	config.maxWait = parseTime(instance.args.maxWait)
	config.randomWait = instance.args.randomWait
	config.minSize = parseSize(instance.args.minSize)
	config.maxSize = parseSize(instance.args.maxSize)
	config.randomSize = instance.args.randomSize
	return config
}

func parseTime(value string) time.Duration {
	parsedDeadTime, err := time.ParseDuration(value)
	if err != nil {
		//This should be A) tested and B) use panic and correctly propogate errors
		fmt.Errorf("cannot parse time from string value")
	}
	return parsedDeadTime
}

func parseSize(value string) uint64 {
	parsedValue, err := humanize.ParseBytes(value)
	if err != nil {
		//This should be A) tested and B) use panic and correctly propogate errors
		fmt.Errorf("cannot parse size from string value")
	}
	return parsedValue
}

func NewArgsConfigurationReader(args CommandLineArgs) *ArgsConfigurationReader {
	return &ArgsConfigurationReader{args, 5 * time.Second}
}

type Configuration struct {
	port       int
	host       string
	verbose    bool
	content    string
	headers    []string
	deadTime   time.Duration
	minWait    time.Duration
	maxWait    time.Duration
	randomWait bool
	minSize    uint64
	maxSize    uint64
	randomSize bool
}
