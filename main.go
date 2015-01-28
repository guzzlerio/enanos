package main

func main() {
	responseBodyGenerator := NewRandomResponseBodyGenerator(10, 10000)
	handleFactory := NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator)
	StartEnanos(responseBodyGenerator, handleFactory)
}
