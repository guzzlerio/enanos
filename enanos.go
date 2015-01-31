package main

import (
	"fmt"
	"net/http"
	"time"
)

type Snoozer interface {
	RandomSnoozeBetween(minDuration time.Duration, max time.Duration)
}

type RealSnoozer struct {
	random Random
}

func (instance *RealSnoozer) RandomSnoozeBetween(min time.Duration, max time.Duration) {
	randomSleep := instance.random.Duration(min, max)
	time.Sleep(randomSleep)
}

func NewRealSnoozer(random Random) *RealSnoozer {
	return &RealSnoozer{random}
}

type Config struct {
	httpHandlerFatory EnanosHttpHandlerFactory
	port              int
	debug             bool
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
	random    Random
}

func (instance *RandomResponseBodyGenerator) Generate() string {
	randValue := instance.random.Int(instance.minLength, instance.maxLength)
	var returnArray = make([]rune, randValue)
	for i := range returnArray {
		returnArray[i] = '-'
	}
	return string(returnArray)
}

func NewRandomResponseBodyGenerator(minLength int, maxLength int) *RandomResponseBodyGenerator {
	random := NewRealRandom()
	return &RandomResponseBodyGenerator{minLength, maxLength, random}
}

type EnanosHttpHandlerFactory interface {
	Happy(w http.ResponseWriter, r *http.Request)
	Grumpy(w http.ResponseWriter, r *http.Request)
	Sneezy(w http.ResponseWriter, r *http.Request)
	Sleepy(w http.ResponseWriter, r *http.Request)
	Bashful(w http.ResponseWriter, r *http.Request)
	Dopey(w http.ResponseWriter, r *http.Request)
}

type DefaultEnanosHttpHandlerFactory struct {
	responseBodyGenerator ResponseBodyGenerator
	snoozer               Snoozer
	random                Random
	responseCodes_300     []int
	responseCodes_400     []int
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

func (instance *DefaultEnanosHttpHandlerFactory) Sleepy(w http.ResponseWriter, r *http.Request) {
	instance.snoozer.RandomSnoozeBetween(1*time.Second, 60*time.Second)
	w.WriteHeader(http.StatusOK)
	data := instance.responseBodyGenerator.Generate()
	w.Write([]byte(data))
}

func (instance *DefaultEnanosHttpHandlerFactory) Bashful(w http.ResponseWriter, r *http.Request) {
	randomIndex := instance.random.Int(0, len(instance.responseCodes_300))
	w.WriteHeader(instance.responseCodes_300[randomIndex])
}

func (instance *DefaultEnanosHttpHandlerFactory) Dopey(w http.ResponseWriter, r *http.Request) {
	randomIndex := instance.random.Int(0, len(instance.responseCodes_400))
	w.WriteHeader(instance.responseCodes_400[randomIndex])
}

func NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator ResponseBodyGenerator, snoozer Snoozer, random Random) *DefaultEnanosHttpHandlerFactory {
	responseCodes_300 := []int{300}
	responseCodes_400 := []int{400}
	return &DefaultEnanosHttpHandlerFactory{responseBodyGenerator, snoozer, random, responseCodes_300, responseCodes_400}
}

func StartEnanos(config Config) {
	mux := http.NewServeMux()
	mux.HandleFunc("/default/happy", func(writer http.ResponseWriter, request *http.Request) {
		if config.debug {
			fmt.Println(fmt.Sprintf("%s - %d bytes - %s", request.RemoteAddr, request.ContentLength, request.URL))
		}
		config.httpHandlerFatory.Happy(writer, request)
	})
	mux.HandleFunc("/default/grumpy", func(writer http.ResponseWriter, request *http.Request) {
		config.httpHandlerFatory.Grumpy(writer, request)
	})
	mux.HandleFunc("/default/sneezy", func(writer http.ResponseWriter, request *http.Request) {
		config.httpHandlerFatory.Sneezy(writer, request)
	})
	mux.HandleFunc("/default/sleepy", func(writer http.ResponseWriter, request *http.Request) {
		config.httpHandlerFatory.Sleepy(writer, request)
	})
	mux.HandleFunc("/default/bashful", func(writer http.ResponseWriter, request *http.Request) {
		config.httpHandlerFatory.Bashful(writer, request)
	})
	mux.HandleFunc("/default/dopey", func(writer http.ResponseWriter, request *http.Request) {
		config.httpHandlerFatory.Dopey(writer, request)
	})
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.port), mux)
	if err != nil {
		fmt.Errorf("error encountered %v", err)
	}
}
