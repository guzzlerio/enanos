package main

import (
	"fmt"
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type ResponseBodyGenerator interface {
	Generate() string
}

type FakeResponseBodyGenerator struct {
	use string
}

func (instance *FakeResponseBodyGenerator) Generate() string {
	return instance.use
}

func (instance *FakeResponseBodyGenerator) UseString(value string) {
	instance.use = value
}

func NewFakeResponseBodyGenerator() *FakeResponseBodyGenerator {
	return &FakeResponseBodyGenerator{""}
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

func Test_Enanos(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Enanos Server:", func() {

		var fakeResponseBodyGenerator *FakeResponseBodyGenerator
		var enanosHttpHandlerFactory *DefaultEnanosHttpHandlerFactory

		url := func(path string) (fullPath string) {
			fullPath = "http://localhost:8000" + path
			return
		}

		g.BeforeEach(func() {
			go func() {
				fakeResponseBodyGenerator = NewFakeResponseBodyGenerator()
				enanosHttpHandlerFactory = NewDefaultEnanosHttpHandlerFactory(fakeResponseBodyGenerator)
				mux := http.NewServeMux()
				mux.HandleFunc("/default/happy", func(writer http.ResponseWriter, request *http.Request) {
					enanosHttpHandlerFactory.Happy(writer, request)
				})
				mux.HandleFunc("/default/grumpy", func(writer http.ResponseWriter, request *http.Request) {
					enanosHttpHandlerFactory.Grumpy(writer, request)
				})
				mux.HandleFunc("/default/sneezy", func(writer http.ResponseWriter, request *http.Request) {
					enanosHttpHandlerFactory.Sneezy(writer, request)
				})
				err := http.ListenAndServe(":8000", mux)
				if err != nil {
					fmt.Errorf("error encountered %v", err)
				}
			}()
		})

		g.Describe("Happy :", func() {
			g.It("GET returns 200", func() {
				resp, _ := http.Get(url("/default/happy"))
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})
		})

		g.Describe("Grumpy :", func() {
			g.It("GET returns 500", func() {
				resp, _ := http.Get(url("/default/grumpy"))
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})
		})

		g.Describe("Sneezy :", func() {
			g.It("GET returns 200", func() {
				resp, _ := http.Get(url("/default/sneezy"))
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})

			g.It("GET returns random response body", func() {
				sample := "foobar"
				fakeResponseBodyGenerator.UseString(sample)
				resp, _ := http.Get(url("/default/sneezy"))
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				assert.Equal(t, sample, string(body))
			})
		})
	})
}
