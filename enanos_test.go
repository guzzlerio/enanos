package main

import (
	"bytes"
	"fmt"
	"github.com/REAANDREW/goclock"
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func ContainsInt(array []int, item int) bool {
	for _, arrayItem := range array {
		if item == arrayItem {
			return true
		}
	}
	return false
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

type FakeSnoozer struct {
	duration time.Duration
}

func (instance *FakeSnoozer) RandomSnoozeBetween(min time.Duration, max time.Duration) {
	time.Sleep(instance.duration)
}

func (instance *FakeSnoozer) SleepFor(duration time.Duration) {
	instance.duration = duration
}

func NewFakeSnoozer() *FakeSnoozer {
	return &FakeSnoozer{0}
}

var (
	fakeResponseBodyGenerator *FakeResponseBodyGenerator
	enanosHttpHandlerFactory  *DefaultEnanosHttpHandlerFactory
	snoozer                   *FakeSnoozer
	random                    *FakeRandom
)

const (
	PORT int = 8000
)

func TestMain(m *testing.M) {
	fakeResponseBodyGenerator = NewFakeResponseBodyGenerator()
	random = NewFakeRandom()
	snoozer = NewFakeSnoozer()
	enanosHttpHandlerFactory = NewDefaultEnanosHttpHandlerFactory(fakeResponseBodyGenerator, snoozer, random)
	go func() {
		config := Config{enanosHttpHandlerFactory, PORT, false}
		StartEnanos(config)
	}()
	os.Exit(m.Run())
}

func Test_ResponseBodyGenerator(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Default Response Body Generator", func() {
		g.It("generates a string of the defined lenth", func() {
			maxLength := 5
			generator := NewDefaultResponseBodyGenerator(maxLength)
			value := generator.Generate()
			assert.Equal(t, maxLength, len(value))
		})
	})

	g.Describe("Random Response Body Generator", func() {
		g.It("generates a string of length between the defined min length and the defined max length", func() {
			minLength := 50
			maxLength := 500
			generator := NewRandomResponseBodyGenerator(minLength, maxLength)
			value := generator.Generate()
			assert.True(t, len(value) >= minLength && len(value) <= maxLength)
		})
	})
}

func SendHelloWorldByHttpMethod(method string, url string) (resp *http.Response, err error) {
	var jsonStr = []byte(`{"message":"hello world"}`)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return
}

func Test_Enanos(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Enanos Server:", func() {

		url := func(path string) (fullPath string) {
			fullPath = fmt.Sprintf("http://localhost:%d", PORT) + path
			return
		}

		g.Describe("Happy :", func() {
			var happyUrl string
			g.Before(func() {
				happyUrl = url("/default/happy")
			})
			g.It("GET returns 200", func() {
				resp, _ := http.Get(happyUrl)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})

			g.It("POST returns 200", func() {
				resp, _ := SendHelloWorldByHttpMethod("POST", happyUrl)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})

			g.It("PUT returns 200", func() {
				resp, _ := SendHelloWorldByHttpMethod("PUT", happyUrl)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})

			g.It("DELETE returns 200", func() {
				resp, _ := SendHelloWorldByHttpMethod("DELETE", happyUrl)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})

			g.It("Any Random Verb returns 200", func() {
				resp, _ := SendHelloWorldByHttpMethod("TALULA", happyUrl)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})
		})

		g.Describe("Grumpy :", func() {
			var grumpyUrl string
			g.Before(func() {
				grumpyUrl = url("/default/grumpy")
			})
			g.It("GET returns 500", func() {
				resp, _ := http.Get(grumpyUrl)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})

			g.It("POST returns 500", func() {
				resp, _ := SendHelloWorldByHttpMethod("POST", grumpyUrl)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})

			g.It("PUT returns 500", func() {
				resp, _ := SendHelloWorldByHttpMethod("PUT", grumpyUrl)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})

			g.It("DELETE returns 500", func() {
				resp, _ := SendHelloWorldByHttpMethod("DELETE", grumpyUrl)
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			})

			g.It("Any Random Verb returns 500", func() {
				resp, _ := SendHelloWorldByHttpMethod("TALULA", grumpyUrl)
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

		g.Describe("Sleepy :", func() {
			g.It("GET returns 200 after a random time between a start and end duration", func() {
				sleep := 10 * time.Millisecond
				snoozer.SleepFor(sleep)
				start := time.Now()
				resp, _ := http.Get(url("/default/sleepy"))
				end := time.Now()
				difference := goclock.DurationDiff(start, end)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				assert.True(t, difference >= sleep && difference <= sleep+(5*time.Millisecond))
			})
		})

		g.Describe("Bashful :", func() {
			g.It("GET returns a 300 response code", func() {
				random.ForIntUse(0)
				resp, _ := http.Get(url("/default/bashful"))
				assert.Equal(t, 300, resp.StatusCode)
			})
			g.It("GET returns a 301 response code")
			g.It("GET returns a 302 response code")
			g.It("GET returns a 304 response code")
			g.It("GET returns a 305 response code")
		})

		g.Describe("Dopey :", func() {
			g.It("GET returns a 400 response code", func() {
				random.ForIntUse(0)
				resp, _ := http.Get(url("/default/dopey"))
				assert.Equal(t, 400, resp.StatusCode)
			})
		})

		g.Describe("Doc :", func() {
			g.It("GET kills the web server and returns after a set time period")
		})
	})
}
