package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	ENV_ENANOS_PORT string = "ENANOS_PORT"
)

var (
	port = kingpin.Flag("port", "the port to host the server on").Default("8000").Short('p').OverrideDefaultFromEnvar(ENV_ENANOS_PORT).Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	responseBodyGenerator := NewRandomResponseBodyGenerator(10, 10000)
	handleFactory := NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator)
	fmt.Println(fmt.Sprintf("Enanos Server listening on port %d", *port))
	StartEnanos(responseBodyGenerator, handleFactory, *port)
}
