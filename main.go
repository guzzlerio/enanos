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
	kingpin.CommandLine.Help = `
	ENANOS

	Is a HTTP server with several endpoints that can be used to substitute the actual http service dependencies of a system.  This tool allows developers to see how the system will perform against varying un-stable http services, each which exhibit different effects. Oh, and the names of the endpoints are the same as those of some well known disney characters:
	Happy   - will return a 200 response code
	Grumpy  - will return a random 5XX response code 
	Sneezy  - will return a 200 response code but a response body with a size between <minSize> and <maxSize>
	Sleepy  - will return a 200 response code but only after a random sleep between <minSleep> and <maxSleep>
	Bashful - will return a random 3XX response code.  If the response code is one which redirects then Bashful will return its own location to invite an infinite redirect loop
	Dopey   - will return a random 4XX response code
	`
	kingpin.Parse()
	responseBodyGenerator := NewRandomResponseBodyGenerator(*minSize, *maxSize)
	random := NewRealRandom()
	snoozer := NewRealSnoozer(time.Duration(*minSleep)*time.Millisecond, time.Duration(*maxSleep)*time.Millisecond)
	handleFactory := NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator, responseCodeGeneratorFactory, snoozer, random)
	config := Config{handleFactory, *port, *debug}
	fmt.Println(fmt.Sprintf("Enanos Server listening on port %d", *port))
	StartEnanos(config)
}
