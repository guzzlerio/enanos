package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	ENV_ENANOS_PORT string = "ENANOS_PORT"
)

var (
	debug = kingpin.Flag("debug", "Enable debug mode.").Bool()
	port  = kingpin.Flag("port", "the port to host the server on").Default("8000").Short('p').OverrideDefaultFromEnvar(ENV_ENANOS_PORT).Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	responseBodyGenerator := NewRandomResponseBodyGenerator(10, 10000)
	random := NewRealRandom()
	snoozer := NewRealSnoozer(random)
	handleFactory := NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator, snoozer)
	config := Config{handleFactory, *port, *debug}
	fmt.Println(fmt.Sprintf("Enanos Server listening on port %d", *port))
	StartEnanos(config)
}
