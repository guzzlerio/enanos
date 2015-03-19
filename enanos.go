package main

import (
	"fmt"
	"github.com/REAANDREW/goSimpleHttp"
	"net/http"
	"sync"
	"time"
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
	host    string
	verbose bool
	content string
	headers []string
}

func StartEnanos(config Config, responseBodyGenerator ResponseBodyGenerator, responseCodeGeneratorFactory func(codes []int) ResponseCodeGenerator, snoozer Snoozer) {
	var shouldStop bool = true
	var wg sync.WaitGroup
	wg.Add(1)
	var handlerFactory HttpHandler = NewDefultHttpHandler(responseBodyGenerator, responseCodeGeneratorFactory, snoozer, config)
	if config.verbose {
		handlerFactory = &VerboseHttpHandler{handlerFactory}
	}
	server := goSimpleHttp.NewSimpleHttpServer(config.port, config.host)
	server.OnStopped(func() {
		if shouldStop {
			wg.Done()
		} else {
			time.Sleep(5 * time.Second)
			server.Start()
			shouldStop = true
		}
	})

	urlToHandlers := map[string]goSimpleHttp.HttpHandler{
		"/success":      handlerFactory.Success,
		"/server_error": handlerFactory.Server_Error,
		"/content_size": handlerFactory.Content_Size,
		"/wait":         handlerFactory.Wait,
		"/redirect":     handlerFactory.Redirect,
		"/client_error": handlerFactory.Client_Error,
		"/defined":      handlerFactory.Defined,
		"/dead_or_alive": func(w http.ResponseWriter, t *http.Request) {
			fmt.Println("Server stopping")
			shouldStop = false
			server.Stop()
		},
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
