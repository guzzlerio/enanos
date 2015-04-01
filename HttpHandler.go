package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HttpResponseWriterRecorder struct {
	Code   int
	writer http.ResponseWriter
}

func (instance *HttpResponseWriterRecorder) Header() http.Header {
	return instance.writer.Header()
}

func (instance *HttpResponseWriterRecorder) Write(data []byte) (int, error) {
	return instance.writer.Write(data)
}

func (instance *HttpResponseWriterRecorder) WriteHeader(code int) {
	instance.Code = code
	instance.writer.WriteHeader(code)
}

func monitorTime(handler http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	writer := &HttpResponseWriterRecorder{0, w}
	handler(writer, r)
	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("%-15s%-4d%s", elapsed, writer.Code, r.URL.Path))
}

type HttpHandler interface {
	Success(w http.ResponseWriter, r *http.Request)
	Server_Error(w http.ResponseWriter, r *http.Request)
	Content_Size(w http.ResponseWriter, r *http.Request)
	Wait(w http.ResponseWriter, r *http.Request)
	Redirect(w http.ResponseWriter, r *http.Request)
	Client_Error(w http.ResponseWriter, r *http.Request)
	Defined(w http.ResponseWriter, r *http.Request)
}

type VerboseHttpHandler struct {
	handler HttpHandler
}

func (instance *VerboseHttpHandler) Success(w http.ResponseWriter, r *http.Request) {
	monitorTime(instance.handler.Success, w, r)
}
func (instance *VerboseHttpHandler) Server_Error(w http.ResponseWriter, r *http.Request) {
	monitorTime(instance.handler.Server_Error, w, r)
}
func (instance *VerboseHttpHandler) Content_Size(w http.ResponseWriter, r *http.Request) {
	monitorTime(instance.handler.Content_Size, w, r)
}
func (instance *VerboseHttpHandler) Wait(w http.ResponseWriter, r *http.Request) {
	monitorTime(instance.handler.Wait, w, r)
}
func (instance *VerboseHttpHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	monitorTime(instance.handler.Redirect, w, r)
}
func (instance *VerboseHttpHandler) Client_Error(w http.ResponseWriter, r *http.Request) {
	monitorTime(instance.handler.Client_Error, w, r)
}
func (instance *VerboseHttpHandler) Defined(w http.ResponseWriter, r *http.Request) {
	monitorTime(instance.handler.Defined, w, r)
}

type DefaultEnanosHttpHandlerFactory struct {
	responseBodyGenerator ResponseBodyGenerator
	responseCodeGenerator ResponseCodeGenerator
	snoozer               Snoozer
	config                Configuration
}

func setHeaders(w http.ResponseWriter, config Configuration) {
	for _, responseHeader := range config.headers {
		split := strings.Split(responseHeader, ":")
		w.Header().Set(split[0], split[1])
	}
}

func (instance *DefaultEnanosHttpHandlerFactory) Success(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, instance.config)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(instance.config.content))
}

func (instance *DefaultEnanosHttpHandlerFactory) Server_Error(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, instance.config)
	code := instance.responseCodeGenerator.GenerateServerErrorCode()
	w.WriteHeader(code)
}

func (instance *DefaultEnanosHttpHandlerFactory) Content_Size(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, instance.config)
	w.WriteHeader(http.StatusOK)
	data := instance.responseBodyGenerator.Generate()
	w.Write([]byte(data))
}

func (instance *DefaultEnanosHttpHandlerFactory) Wait(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, instance.config)
	instance.snoozer.Snooze()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(instance.config.content))
}

func (instance *DefaultEnanosHttpHandlerFactory) Redirect(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, instance.config)
	code := instance.responseCodeGenerator.GenerateRedirectionCode()
	if code == 301 || code == 302 || code == 303 || code == 307 {
		existingHeader := w.Header().Get("location")
		if existingHeader == "" {
			w.Header().Set("location", "/redirect")
		}
	}
	w.WriteHeader(code)
}

func (instance *DefaultEnanosHttpHandlerFactory) Client_Error(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, instance.config)
	code := instance.responseCodeGenerator.GenerateClientErrorCode()
	w.WriteHeader(code)
}

func (instance *DefaultEnanosHttpHandlerFactory) Defined(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, instance.config)
	code := r.URL.Query().Get("code")
	if code != "" {
		intCode, err := strconv.Atoi(code)
		if err != nil {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(intCode)
		}
	} else {
		w.WriteHeader(400)
	}
}

func NewDefultHttpHandler(responseBodyGenerator ResponseBodyGenerator, responseCodeGenerator ResponseCodeGenerator, snoozer Snoozer, config Configuration) *DefaultEnanosHttpHandlerFactory {
	return &DefaultEnanosHttpHandlerFactory{responseBodyGenerator, responseCodeGenerator, snoozer, config}
}
