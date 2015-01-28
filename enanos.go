package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

type ResponseBodyGenerator interface {
	Generate() string
}

type DefaultResponseBodyGenerator struct {
	maxLength int
}

func (instance *DefaultResponseBodyGenerator) Generate() string {
	var returnArray = make([]rune, instance.maxLength)
	for i := range returnArray {
		returnArray[i] = '-'
	}
	return string(returnArray)
}

func NewDefaultResponseBodyGenerator(maxLength int) *DefaultResponseBodyGenerator {
	return &DefaultResponseBodyGenerator{maxLength}
}

type RandomResponseBodyGenerator struct {
	minLength int
	maxLength int
}

func (instance *RandomResponseBodyGenerator) Generate() string {
	var returnArray = make([]rune, random(instance.minLength, instance.maxLength))
	for i := range returnArray {
		returnArray[i] = '-'
	}
	return string(returnArray)
}

func NewRandomResponseBodyGenerator(minLength int, maxLength int) *RandomResponseBodyGenerator {
	return &RandomResponseBodyGenerator{minLength, maxLength}
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
