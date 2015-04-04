package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type CommandLineArgs struct {
	Port       int
	Host       string
	Verbose    bool
	Content    string
	DeadTime   string
	MinWait    string
	MaxWait    string
	RandomWait bool
	MinSize    string
	MaxSize    string
	RandomSize bool
	Config     string
	Headers    []string
}

type ConfigurationReader interface {
	Read() Configuration
}

type ArgsConfigurationReader struct {
	args        *CommandLineArgs
	defaultTime time.Duration
}

func (instance *ArgsConfigurationReader) Read() Configuration {
	config := Configuration{}
	if instance.args.Config != "" {
		data, err := ioutil.ReadFile(instance.args.Config)
		if err != nil {
			fmt.Errorf("Cannot read the path for the config file")
		}
		err = yaml.Unmarshal(data, instance.args)
		if err != nil {
			fmt.Errorf("Cannot read the config yml")
		}
	}
	config.port = instance.args.Port
	config.host = instance.args.Host
	config.verbose = instance.args.Verbose
	config.content = instance.args.Content
	config.headers = instance.args.Headers
	config.deadTime = parseTime(instance.args.DeadTime)
	config.minWait = parseTime(instance.args.MinWait)
	config.maxWait = parseTime(instance.args.MaxWait)
	config.randomWait = instance.args.RandomWait
	config.minSize = parseSize(instance.args.MinSize)
	config.maxSize = parseSize(instance.args.MaxSize)
	config.randomSize = instance.args.RandomSize
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

func NewArgsConfigurationReader(args *CommandLineArgs) *ArgsConfigurationReader {
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
