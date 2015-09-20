package main

import (
	"github.com/reaandrew/goSimpleHttp"
	"net/http"
	"sync"
	"time"
    "fmt"
)

const (
	STOPPED_EVENT_KEY string = "stopped"
)

var (
	responseCodes_300 []int = []int{300, 301, 302, 303, 304, 305, 307}
	responseCodes_400 []int = []int{400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 429}
	responseCodes_500 []int = []int{500, 501, 502, 503, 504, 505}
)

type Server interface{
    Start()
    Stop()
}

type JitterServer struct{
    Config Configuration
    ResponseBodyGenerator ResponseBodyGenerator
    ResponseCodeGenerator ResponseCodeGenerator
    Snoozer Snoozer
    Server *goSimpleHttp.SimpleHttpServer
}

func (instance *JitterServer) Start(){
    config := instance.Config
	instance.Server = goSimpleHttp.NewSimpleHttpServer(config.port + 1, config.host)
    if config.jitterTime == time.Duration(0){
        return
    }
	var handlerFactory HttpHandler = NewDefultHttpHandler(instance.ResponseBodyGenerator, instance.ResponseCodeGenerator, instance.Snoozer, instance.Config)
	if config.verbose {
		handlerFactory = &VerboseHttpHandler{handlerFactory}
	}
    ticker := time.NewTicker(config.jitterTime)
    stopped := false
    go func() {
        for {
           select {
            case <- ticker.C:
                if stopped {
                    fmt.Println("Starting server")
                    instance.Server.Start()
                    stopped = false
                }else{
                    fmt.Println("Stopping server")
                    instance.Server.Stop()
                    stopped = true
                }
            }
        }
     }()
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
		instance.Server.Get(key, value)
		instance.Server.Post(key, value)
		instance.Server.Put(key, value)
		instance.Server.Delete(key, value)
	}

	instance.Server.Start()
}
func (instance *JitterServer) Stop(){
    instance.Server.Stop()
}

type HarnessServer struct{
    Config Configuration
    ResponseBodyGenerator ResponseBodyGenerator
    ResponseCodeGenerator ResponseCodeGenerator
    Snoozer Snoozer
    Server *goSimpleHttp.SimpleHttpServer
}

func (instance *HarnessServer) Start(){
    config := instance.Config
	var shouldStop bool = true
	var handlerFactory HttpHandler = NewDefultHttpHandler(instance.ResponseBodyGenerator, instance.ResponseCodeGenerator, instance.Snoozer, instance.Config)
	if config.verbose {
		handlerFactory = &VerboseHttpHandler{handlerFactory}
	}
	instance.Server = goSimpleHttp.NewSimpleHttpServer(config.port, config.host)
	instance.Server.OnStopped(func() {
        time.Sleep(config.deadTime)
        instance.Server.Start()
        shouldStop = true
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
			shouldStop = false
			instance.Server.Stop()
		},
	}

	for key, value := range urlToHandlers {
		instance.Server.Get(key, value)
		instance.Server.Post(key, value)
		instance.Server.Put(key, value)
		instance.Server.Delete(key, value)
	}

	instance.Server.Start()
}

func (instance *HarnessServer) Stop(){
    instance.Server.Stop()
}

type EnanosServer struct{
    Servers []Server
    WaitHandle *sync.WaitGroup
}

func (instance *EnanosServer) Start(){
    for _,server := range instance.Servers{
        server.Start()
    }
	 instance.WaitHandle.Add(1)
}

func (instance *EnanosServer) Stop(){
    for _,server := range instance.Servers{
        server.Stop()
    }
	 fmt.Println("Calling Done")
    instance.WaitHandle.Done()
}

type ServerFactory struct{
    Config Configuration
    ResponseBodyGenerator ResponseBodyGenerator
    ResponseCodeGenerator ResponseCodeGenerator
    Snoozer Snoozer
    WaitHandle *sync.WaitGroup
}

func (instance *ServerFactory) CreateServer() Server{
    jitterServer := &JitterServer{
        Config : instance.Config,
        ResponseBodyGenerator : instance.ResponseBodyGenerator,
        ResponseCodeGenerator : instance.ResponseCodeGenerator,
        Snoozer : instance.Snoozer,
    }

    harnessServer := &HarnessServer{
        Config : instance.Config,
        ResponseBodyGenerator : instance.ResponseBodyGenerator,
        ResponseCodeGenerator : instance.ResponseCodeGenerator,
        Snoozer : instance.Snoozer,
    }

    servers := []Server{jitterServer, harnessServer}

    return &EnanosServer{
        Servers: servers,
        WaitHandle : instance.WaitHandle,
    }
}
