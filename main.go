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
	kingpin.CommandLine.Help = `Enanos is a HTTP server with several endpoints that can be used to substitute the actual http service dependencies of a system.  This tool allows developers to see how a system will perform against varying un-stable http services, each which exhibit different effects.
	
	/success		- will return a 200 response code
	/server_error		- will return a random 5XX response code 
	/content_size		- will return a 200 response code but a response body with a size between <minSize> and <maxSize>.  The content returned will be random or a mangled version of the content which has been configured to return i.e. it cannot guarantee to meet any content-types configured in that it will be malformed.
	/wait			- will return a 200 response code but only after a random sleep between <minSleep> and <maxSleep>
	/redirect		- will return a random 3XX response code.  If the response code is one which redirects then Bashful will return its own location to invite an infinite redirect loop
	/client_error		- will return a random 4XX response code
	/dead_or_alive	- will kill the server and only bring it back online after configured amount of time (ms) has passed

	/defined?code=<code>	- will return the specified http status code
	`
	kingpin.Parse()
	responseBodyGenerator := NewRandomResponseBodyGenerator(*minSize, *maxSize)
	snoozer := NewRealSnoozer(time.Duration(*minSleep)*time.Millisecond, time.Duration(*maxSleep)*time.Millisecond)

	config := Config{*port, *debug, "sample", "application/xml"}
	fmt.Println(fmt.Sprintf("Enanos Server listening on port %d", *port))
	StartEnanos(config, responseBodyGenerator, responseCodeGeneratorFactory, snoozer)
}
