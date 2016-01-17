package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.WriteHeader(http.StatusCreated)
	} else {
		io.WriteString(w, "Hello world!")
	}
}

//HTTPServer ...
type HTTPServer struct {
	Port     int
	listener net.Listener
	server   *http.Server
	mux      *http.ServeMux
}

//NewHTTPServer ...
func NewHTTPServer(port int) *HTTPServer {
	return &HTTPServer{
		Port: port,
	}
}

//Start ...
func (instance *HTTPServer) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", instance.Port))

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err != nil {
		return err
	}
	instance.listener = l
	instance.mux = mux
	instance.server = s

	go func(listener net.Listener) {
		s.Serve(listener)
	}(l)

	return nil
}

//Stop ...
func (instance *HTTPServer) Stop() {
	if instance.listener != nil {
		instance.listener.Close()
	}
}
