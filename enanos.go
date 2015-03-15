package main

import (
	"github.com/REAANDREW/goSimpleHttp"
	"sync"
)

const (
	STOPPED_EVENT_KEY string = "stopped"
)

var (
	responseCodes_300 []int = []int{300, 301, 302, 303, 304, 305, 307}
	responseCodes_400 []int = []int{400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417}
	responseCodes_500 []int = []int{500, 501, 502, 503, 504, 505}
)

type Config struct {
	port    int
	debug   bool
	content string
}

func StartEnanos(config Config, responseBodyGenerator ResponseBodyGenerator, responseCodeGeneratorFactory func(codes []int) ResponseCodeGenerator, snoozer Snoozer) {
	var wg sync.WaitGroup
	wg.Add(1)
	handlerFactory := NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator, responseCodeGeneratorFactory, snoozer, config)
	server := goSimpleHttp.NewSimpleHttpServer(config.port, "localhost")
	server.OnStopped(func() {
		wg.Done()
	})

	urlToHandlers := map[string]goSimpleHttp.HttpHandler{
		"/success":      handlerFactory.Success,
		"/server_error": handlerFactory.Server_Error,
		"/content_size": handlerFactory.Content_Size,
		"/wait":         handlerFactory.Wait,
		"/redirect":     handlerFactory.Redirect,
		"/client_error": handlerFactory.Client_Error,
		"/defined":      handlerFactory.Defined,
	}

	for key, value := range urlToHandlers {
		server.Get(key, value)
		server.Post(key, value)
		server.Put(key, value)
		server.Delete(key, value)
	}

	server.Start()
	wg.Wait()
}
