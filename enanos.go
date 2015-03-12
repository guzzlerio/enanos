package main

import (
	"github.com/REAANDREW/goSimpleHttp"
	"net/http"
	"sync"
	"time"
)

const (
	STOPPED_EVENT_KEY string = "stopped"
)

type ResponseCodeGenerator interface {
	Generate() int
}

type RandomResponseCodeGenerator struct {
	responseCodes []int
	randomGen     Random
}

func (instance *RandomResponseCodeGenerator) Generate() int {
	from := 0
	to := len(instance.responseCodes)
	index := instance.randomGen.Int(from, to)
	return instance.responseCodes[index]
}

func NewRandomResponseCodeGenerator(responseCodes []int) *RandomResponseCodeGenerator {
	return &RandomResponseCodeGenerator{responseCodes, NewRealRandom()}
}

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
	responseCodes_300     ResponseCodeGenerator
	responseCodes_400     ResponseCodeGenerator
	responseCodes_500     ResponseCodeGenerator
}

func (instance *DefaultEnanosHttpHandlerFactory) Happy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (instance *DefaultEnanosHttpHandlerFactory) Grumpy(w http.ResponseWriter, r *http.Request) {
	code := instance.responseCodes_500.Generate()
	w.WriteHeader(code)
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
	code := instance.responseCodes_300.Generate()
	w.WriteHeader(code)
}

func (instance *DefaultEnanosHttpHandlerFactory) Dopey(w http.ResponseWriter, r *http.Request) {
	code := instance.responseCodes_400.Generate()
	w.WriteHeader(code)
}

func NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator ResponseBodyGenerator, responseCodeGenFactory func(codes []int) ResponseCodeGenerator, snoozer Snoozer, random Random) *DefaultEnanosHttpHandlerFactory {
	responseCodes_300 := responseCodeGenFactory([]int{300})
	responseCodes_400 := responseCodeGenFactory([]int{400})
	responseCodes_500 := responseCodeGenFactory([]int{500})
	return &DefaultEnanosHttpHandlerFactory{responseBodyGenerator, snoozer, random, responseCodes_300, responseCodes_400, responseCodes_500}
}

func StartEnanos(config Config) {
	var wg sync.WaitGroup
	wg.Add(1)
	server := goSimpleHttp.NewSimpleHttpServer(config.port, "localhost")
	server.OnStopped(func() {
		wg.Done()
	})

	urlToHandlers := map[string]goSimpleHttp.HttpHandler{
		"/default/happy":   config.httpHandlerFatory.Happy,
		"/default/grumpy":  config.httpHandlerFatory.Grumpy,
		"/default/sneezy":  config.httpHandlerFatory.Sneezy,
		"/default/sleepy":  config.httpHandlerFatory.Sleepy,
		"/default/bashful": config.httpHandlerFatory.Bashful,
		"/default/dopey":   config.httpHandlerFatory.Dopey,
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
