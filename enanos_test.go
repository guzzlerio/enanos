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
	"strconv"
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

var (
	fakeResponseBodyGenerator *FakeResponseBodyGenerator
	enanosHttpHandlerFactory  *DefaultEnanosHttpHandlerFactory
	snoozer                   *FakeSnoozer
	responseCodeGenerator     *FakeResponseCodeGenerator
	METHODS                   []string = []string{"GET", "POST", "PUT", "DELETE"}
	content                   string
	contentType               string
)

func factory(codes []int) ResponseCodeGenerator {
	return responseCodeGenerator
}

const (
	PORT int = 8000
)

func TestMain(m *testing.M) {
	fakeResponseBodyGenerator = NewFakeResponseBodyGenerator()
	snoozer = NewFakeSnoozer()
	responseCodeGenerator = NewFakeResponseCodeGenerator()
	content = "<xml type=\"foobar\"></xml>"
	contentType = "application/json"
	go func() {
		config := Config{PORT, false, content, contentType}
		StartEnanos(config, fakeResponseBodyGenerator, factory, snoozer)
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
	return
}

func Test_Enanos(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Enanos Server:", func() {

		url := func(path string) (fullPath string) {
			fullPath = fmt.Sprintf("http://localhost:%d", PORT) + path
			return
		}

		g.Describe("Success :", func() {
			for _, method := range METHODS {
				g.Describe(fmt.Sprintf("%s :", method), func() {
					g.It(fmt.Sprintf("%s returns 200", method), func() {
						resp, _ := SendHelloWorldByHttpMethod(method, url("/success"))
						defer resp.Body.Close()
						assert.Equal(t, http.StatusOK, resp.StatusCode)
					})
				})
			}

			g.It("Returns defined content", func() {
				resp, _ := SendHelloWorldByHttpMethod("GET", url("/success"))
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				assert.Equal(t, string(body), content)
			})

			g.It("Returns defined content-type", func() {
				resp, _ := SendHelloWorldByHttpMethod("GET", url("/success"))
				defer resp.Body.Close()
				contentType := resp.Header.Get("content-type")
				assert.Equal(t, contentType, contentType)
			})
		})

		g.Describe("Server Error :", func() {
			codes := responseCodes_500
			for _, method := range METHODS {
				g.Describe(fmt.Sprintf("%s :", method), func() {
					for _, code := range codes {
						g.It(fmt.Sprintf("%s returns a %d response code", method, code), func() {
							responseCodeGenerator.Use(code)
							resp, _ := SendHelloWorldByHttpMethod(method, url("/server_error"))
							defer resp.Body.Close()
							assert.Equal(t, code, resp.StatusCode)
						})
					}
				})
			}
		})

		g.Describe("Content Size :", func() {
			for _, method := range METHODS {
				g.Describe(fmt.Sprintf("%s :", method), func() {
					g.It(fmt.Sprintf("%s returns 200", method), func() {
						resp, _ := SendHelloWorldByHttpMethod(method, url("/content_size"))
						defer resp.Body.Close()
						assert.Equal(t, http.StatusOK, resp.StatusCode)
					})
				})
				g.It(fmt.Sprintf("%s returns random response body", method), func() {
					sample := "foobar"
					fakeResponseBodyGenerator.UseString(sample)
					resp, _ := SendHelloWorldByHttpMethod(method, url("/content_size"))
					defer resp.Body.Close()
					body, _ := ioutil.ReadAll(resp.Body)
					assert.Equal(t, sample, string(body))
				})
			}
		})

		g.Describe("Wait :", func() {
			for _, method := range METHODS {
				g.Describe(fmt.Sprintf("%s :", method), func() {
					g.It(fmt.Sprintf("%s returns 200 after a random time between a start and end duration", method), func() {
						sleep := 10 * time.Millisecond
						snoozer.SleepFor(sleep)
						start := time.Now()
						resp, _ := SendHelloWorldByHttpMethod(method, url("/wait"))
						defer resp.Body.Close()
						end := time.Now()
						difference := goclock.DurationDiff(start, end)
						assert.Equal(t, http.StatusOK, resp.StatusCode)

						assert.True(t, difference >= sleep && difference <= sleep+(5*time.Millisecond))
					})
				})
			}
		})

		g.Describe("Redirect :", func() {
			codes := responseCodes_300
			for _, method := range METHODS {
				g.Describe(fmt.Sprintf("%s :", method), func() {
					for _, code := range codes {
						g.It(fmt.Sprintf("%s returns a %d response code", method, code), func() {
							responseCodeGenerator.Use(code)
							resp, _ := SendHelloWorldByHttpMethod(method, url("/redirect"))
							defer resp.Body.Close()
							assert.Equal(t, code, resp.StatusCode)
						})
					}
				})
			}
		})

		g.Describe("Client Error :", func() {
			codes := responseCodes_400
			for _, method := range METHODS {
				g.Describe(fmt.Sprintf("%s :", method), func() {
					for _, code := range codes {
						g.It(fmt.Sprintf("%s returns a %d response code", method, code), func() {
							responseCodeGenerator.Use(code)
							resp, _ := SendHelloWorldByHttpMethod(method, url("/client_error"))
							defer resp.Body.Close()
							assert.Equal(t, code, resp.StatusCode)
						})
					}
				})
			}
		})

		g.Describe("Doc :", func() {
			g.It("GET kills the web server and returns after a set time period")
		})

		g.Describe("Defined", func() {
			codes := append(responseCodes_300, responseCodes_400...)
			codes = append(codes, responseCodes_500...)
			for _, method := range METHODS {
				g.Describe(fmt.Sprintf("%s :", method), func() {
					for _, code := range codes {
						g.It(fmt.Sprintf("%s returns a %d response code", method, code), func() {
							resp, _ := SendHelloWorldByHttpMethod(method, url("/defined?code="+strconv.Itoa(code)))
							defer resp.Body.Close()
							assert.Equal(t, code, resp.StatusCode)
						})
					}
				})
			}

			g.It("returns 400 when no code is present", func() {
				for _, method := range METHODS {
					code := 400
					resp, _ := SendHelloWorldByHttpMethod(method, url("/defined"))
					defer resp.Body.Close()
					assert.Equal(t, code, resp.StatusCode)
				}
			})

			g.It("returns 400 when an non int code is specified", func() {
				for _, method := range METHODS {
					code := 400
					resp, _ := SendHelloWorldByHttpMethod(method, url("/defined?code=bang"))
					defer resp.Body.Close()
					assert.Equal(t, code, resp.StatusCode)
				}
			})
		})
	})
}
