package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v1"
	"strconv"
	"time"
)

const (
	ENV_ENANOS_PORT string = "ENANOS_PORT"
)

var (
	debug    = kingpin.Flag("debug", "Enable debug mode.").Bool()
	port     = kingpin.Flag("port", "the port to host the server on").Default("8000").Short('p').OverrideDefaultFromEnvar(ENV_ENANOS_PORT).Int()
	minSleep = kingpin.Flag("min-sleep", "the minimum sleep time for sleepy in milliseconds").Default("1000").Int()
	maxSleep = kingpin.Flag("max-sleep", "the maximum sleep time for sleepy in milliseconds").Default("60000").Int()
	minSize  = kingpin.Flag("min-size", "the minimum size of response body for sneezy to generate").Default("1024").Int()
	maxSize  = kingpin.Flag("max-size", "the maximum size of response body for sneezy to generate").Default(strconv.Itoa(1024 * 100)).Int()
)

func responseCodeGeneratorFactory(codes []int) ResponseCodeGenerator {
	return NewRandomResponseCodeGenerator(codes)
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	responseBodyGenerator := NewRandomResponseBodyGenerator(*minSize, *maxSize)
	random := NewRealRandom()
	snoozer := NewRealSnoozer(time.Duration(*minSleep)*time.Millisecond, time.Duration(*maxSleep)*time.Millisecond)
	handleFactory := NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator, responseCodeGeneratorFactory, snoozer, random)
	config := Config{handleFactory, *port, *debug}
	fmt.Println(fmt.Sprintf("Enanos Server listening on port %d", *port))
	StartEnanos(config)
}
