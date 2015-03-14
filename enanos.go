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

var (
	responseCodes_300 []int = []int{300, 301, 302, 303, 304, 305, 307}
	responseCodes_400 []int = []int{400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417}
	responseCodes_500 []int = []int{500, 501, 502, 503, 504, 505}
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
	RandomSnooze()
}

type RealSnoozer struct {
	Min    time.Duration
	Max    time.Duration
	random Random
}

func (instance *RealSnoozer) RandomSnooze() {
	randomSleep := instance.random.Duration(instance.Min, instance.Max)
	time.Sleep(randomSleep)
}

func NewRealSnoozer(min time.Duration, max time.Duration) *RealSnoozer {
	return &RealSnoozer{min, max, &RealRandom{}}
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
	instance.snoozer.RandomSnooze()
	w.WriteHeader(http.StatusOK)
	data := instance.responseBodyGenerator.Generate()
	w.Write([]byte(data))
}

func (instance *DefaultEnanosHttpHandlerFactory) Bashful(w http.ResponseWriter, r *http.Request) {
	code := instance.responseCodes_300.Generate()
	if code == 301 || code == 302 || code == 303 || code == 307 {
		w.Header().Set("location", "/default/bashful")
	}
	w.WriteHeader(code)
}

func (instance *DefaultEnanosHttpHandlerFactory) Dopey(w http.ResponseWriter, r *http.Request) {
	code := instance.responseCodes_400.Generate()
	w.WriteHeader(code)
}

func NewDefaultEnanosHttpHandlerFactory(responseBodyGenerator ResponseBodyGenerator, responseCodeGenFactory func(codes []int) ResponseCodeGenerator, snoozer Snoozer, random Random) *DefaultEnanosHttpHandlerFactory {
	responseCodes_300 := responseCodeGenFactory(responseCodes_300)
	responseCodes_400 := responseCodeGenFactory(responseCodes_400)
	responseCodes_500 := responseCodeGenFactory(responseCodes_500)
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
		"/success":      config.httpHandlerFatory.Happy,
		"/server_error": config.httpHandlerFatory.Grumpy,
		"/content_size": config.httpHandlerFatory.Sneezy,
		"/wait":         config.httpHandlerFatory.Sleepy,
		"/redirect":     config.httpHandlerFatory.Bashful,
		"/client_error": config.httpHandlerFatory.Dopey,
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
