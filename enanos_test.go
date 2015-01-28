package main

import (
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

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

func Test_Enanos_Without_Goblin(t *testing.T) {
	var fakeResponseBodyGenerator *FakeResponseBodyGenerator = NewFakeResponseBodyGenerator()
	var enanosHttpHandlerFactory *DefaultEnanosHttpHandlerFactory = NewDefaultEnanosHttpHandlerFactory(fakeResponseBodyGenerator)

	url := func(path string) (fullPath string) {
		fullPath = "http://localhost:8000" + path
		return
	}
	go func() {
		StartEnanos(fakeResponseBodyGenerator, enanosHttpHandlerFactory)
	}()
	sample := "foobar"
	fakeResponseBodyGenerator.UseString(sample)
	resp, _ := http.Get(url("/default/sneezy"))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, sample, string(body))
}

func Test_DefaultResponseBodyGenerator(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Default Response Body Generator", func() {
		g.It("generates a string of the defined lenth", func() {
			maxLength := 5
			generator := NewDefaultResponseBodyGenerator(maxLength)
			value := generator.Generate()
			assert.Equal(t, maxLength, len(value))
		})
	})
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
			fakeResponseBodyGenerator = NewFakeResponseBodyGenerator()
			enanosHttpHandlerFactory = NewDefaultEnanosHttpHandlerFactory(fakeResponseBodyGenerator)
			go func() {
				StartEnanos(fakeResponseBodyGenerator, enanosHttpHandlerFactory)
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
