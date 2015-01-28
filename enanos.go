package main

import (
	"fmt"
	"net/http"
)

type ResponseBodyGenerator interface {
	Generate() string
}

type EnanosHttpHandlerFactory interface {
	Happy(w http.ResponseWriter, r *http.Request)
	Grumpy(w http.ResponseWriter, r *http.Request)
	Sneezy(w http.ResponseWriter, r *http.Request)
}

type DefaultEnanosHttpHandlerFactory struct {
	responseBodyGenerator ResponseBodyGenerator
}

func (instance *DefaultEnanosHttpHandlerFactory) Happy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (instance *DefaultEnanosHttpHandlerFactory) Grumpy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (instance *DefaultEnanosHttpHandlerFactory) Sneezy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	data := instance.responseBodyGenerator.Generate()
	w.Write([]byte(data))
}

func NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator ResponseBodyGenerator) *DefaultEnanosHttpHandlerFactory {
	return &DefaultEnanosHttpHandlerFactory{responseBodyGenerator}
}

func StartEnanos(responseBodyGenerator ResponseBodyGenerator, handlerFactory EnanosHttpHandlerFactory) {
	mux := http.NewServeMux()
	mux.HandleFunc("/default/happy", func(writer http.ResponseWriter, request *http.Request) {
		handlerFactory.Happy(writer, request)
	})
	mux.HandleFunc("/default/grumpy", func(writer http.ResponseWriter, request *http.Request) {
		handlerFactory.Grumpy(writer, request)
	})
	mux.HandleFunc("/default/sneezy", func(writer http.ResponseWriter, request *http.Request) {
		handlerFactory.Sneezy(writer, request)
	})
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Errorf("error encountered %v", err)
	}
}
