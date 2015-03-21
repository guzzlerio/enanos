package main

import (
	"bytes"
	"fmt"
	"github.com/REAANDREW/goclock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	fakeResponseBodyGenerator *FakeResponseBodyGenerator
	enanosHttpHandlerFactory  *DefaultEnanosHttpHandlerFactory
	snoozer                   *FakeSnoozer
	responseCodeGenerator     *FakeResponseCodeGenerator
	METHODS                   []string = []string{"GET", "POST", "PUT", "DELETE"}
	testContent               string
	testContentType           string
	testHeaders               []string
)

const (
	PORT int = 8000
)

func TestMain(m *testing.M) {
	fakeResponseBodyGenerator = NewFakeResponseBodyGenerator()
	snoozer = NewFakeSnoozer()
	responseCodeGenerator = NewFakeResponseCodeGenerator()
	testContent = "<xml type=\"foobar\"></xml>"
	testContentType = "application/json"
	testHeaders = []string{
		"Age:12",
		"Content-Length:101",
		"Content-Type:" + testContentType,
	}
	go func() {
		config := Config{PORT, "localhost", false, testContent, testHeaders}
		StartEnanos(config, fakeResponseBodyGenerator, responseCodeGenerator, snoozer)
	}()
	os.Exit(m.Run())
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

var _ = Describe("Enanos Server:", func() {

	url := func(path string) (fullPath string) {
		fullPath = fmt.Sprintf("http://localhost:%d", PORT) + path
		return
	}

	Describe("Success :", func() {
		for _, method := range METHODS {
			Describe(fmt.Sprintf("%s :", method), func() {
				It(fmt.Sprintf("%s returns 200", method), func() {
					resp, _ := SendHelloWorldByHttpMethod(method, url("/success"))
					defer resp.Body.Close()
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})
			})
		}

		It("Returns defined content", func() {
			resp, _ := SendHelloWorldByHttpMethod("GET", url("/success"))
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			Expect(string(body)).To(Equal(testContent))
		})

		It("Returns defined content-type", func() {
			resp, _ := SendHelloWorldByHttpMethod("GET", url("/success"))
			defer resp.Body.Close()
			contentType := resp.Header.Get("content-type")
			Expect(contentType).To(Equal(testContentType))
		})
	})

	Describe("Server Error :", func() {
		codes := responseCodes_500
		for _, method := range METHODS {
			Describe(fmt.Sprintf("%s :", method), func() {
				for _, code := range codes {
					It(fmt.Sprintf("%s returns a %d response code", method, code), func() {
						responseCodeGenerator.Use(code)
						resp, _ := SendHelloWorldByHttpMethod(method, url("/server_error"))
						defer resp.Body.Close()
						Expect(code).To(Equal(resp.StatusCode))
					})
				}
			})
		}
	})

	Describe("Content Size :", func() {
		for _, method := range METHODS {
			Describe(fmt.Sprintf("%s :", method), func() {
				It(fmt.Sprintf("%s returns 200", method), func() {
					resp, _ := SendHelloWorldByHttpMethod(method, url("/content_size"))
					defer resp.Body.Close()
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})
			})
			It(fmt.Sprintf("%s returns random response body", method), func() {
				sample := "foobar"
				fakeResponseBodyGenerator.UseString(sample)
				resp, _ := SendHelloWorldByHttpMethod(method, url("/content_size"))
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				Expect(string(body)).To(Equal(sample))
			})
		}
	})

	Describe("Wait :", func() {
		for _, method := range METHODS {
			Describe(fmt.Sprintf("%s :", method), func() {
				It(fmt.Sprintf("%s returns 200 after a random time between a start and end duration", method), func() {
					sleep := 10 * time.Millisecond
					snoozer.SleepFor(sleep)
					start := time.Now()
					resp, _ := SendHelloWorldByHttpMethod(method, url("/wait"))
					defer resp.Body.Close()
					end := time.Now()
					difference := goclock.DurationDiff(start, end)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(difference >= sleep && difference <= sleep+(5*time.Millisecond)).To(BeTrue())
				})
			})
		}
	})

	Describe("Redirect :", func() {
		codes := responseCodes_300
		for _, method := range METHODS {
			Describe(fmt.Sprintf("%s :", method), func() {
				for _, code := range codes {
					It(fmt.Sprintf("%s returns a %d response code", method, code), func() {
						responseCodeGenerator.Use(code)
						resp, _ := SendHelloWorldByHttpMethod(method, url("/redirect"))
						defer resp.Body.Close()
						Expect(resp.StatusCode).To(Equal(code))
						Expect(resp.Header.Get("location")).To(Equal("/redirect"))
					})
				}
			})
		}
	})

	Describe("Client Error :", func() {
		codes := responseCodes_400
		for _, method := range METHODS {
			Describe(fmt.Sprintf("%s :", method), func() {
				for _, code := range codes {
					It(fmt.Sprintf("%s returns a %d response code", method, code), func() {
						responseCodeGenerator.Use(code)
						resp, _ := SendHelloWorldByHttpMethod(method, url("/client_error"))
						defer resp.Body.Close()
						Expect(resp.StatusCode).To(Equal(code))
					})
				}
			})
		}
	})

	Describe("Doc :", func() {
		It("GET kills the web server and returns after a set time period", func() {})
	})

	Describe("Defined", func() {
		codes := append(responseCodes_300, responseCodes_400...)
		codes = append(codes, responseCodes_500...)
		for _, method := range METHODS {
			Describe(fmt.Sprintf("%s :", method), func() {
				for _, code := range codes {
					It(fmt.Sprintf("%s returns a %d response code", method, code), func() {
						resp, _ := SendHelloWorldByHttpMethod(method, url("/defined?code="+strconv.Itoa(code)))
						defer resp.Body.Close()
						Expect(resp.StatusCode).To(Equal(code))
					})
				}
			})
		}

		It("returns 400 when no code is present", func() {
			for _, method := range METHODS {
				code := 400
				resp, _ := SendHelloWorldByHttpMethod(method, url("/defined"))
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(code))
			}
		})

		It("returns 400 when an non int code is specified", func() {
			for _, method := range METHODS {
				code := 400
				resp, _ := SendHelloWorldByHttpMethod(method, url("/defined?code=bang"))
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(code))
			}
		})
	})

	Describe("Headers", func() {
		endpoints := []string{"success", "wait", "content_size"}
		BeforeEach(func() {
			sample := "foobar"
			fakeResponseBodyGenerator.UseString(sample)
			sleep := 10 * time.Millisecond
			snoozer.SleepFor(sleep)
		})

		for _, endpoint := range endpoints {
			for _, method := range METHODS {
				Describe(fmt.Sprintf("%s :", method), func() {
					It(fmt.Sprintf("%s %s set the response headers", method, endpoint), func() {
						urlToUse := url(fmt.Sprintf("/%s", endpoint))
						resp, _ := SendHelloWorldByHttpMethod(method, urlToUse)
						defer resp.Body.Close()

						Expect(resp.Header.Get("Age")).To(Equal("12"))
						Expect(resp.Header.Get("Content-Length")).To(Equal("101"))
					})
				})
			}
		}
	})

})
