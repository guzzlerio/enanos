package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var _ = Describe("Configuration", func() {

	Describe("FileConfigurationReader", func() {

		Describe("reads", func() {

			var config Configuration
			var file *os.File
			var err error
			var headers []string

			BeforeEach(func() {
				headers = []string{"Age:1", "Content-type:text/plain"}
				var data = `port: 8080
host: foobar
verbose: true
content: boom
deadtime: 10s
minwait: 11s
maxwait: 12s
randomwait: true
minsize: 1KB
maxsize: 2KB
randomsize: true
headers: ["Age:1","Content-type:text/plain"]`
				file, err = ioutil.TempFile("", "enanos")
				file.WriteString(data)
				file.Close()
				check(err)

				args := CommandLineArgs{}
				args.Config = file.Name()

				config = NewArgsConfigurationReader(&args).Read()
			})

			AfterEach(func() {
				os.Remove(file.Name())
			})

			It("port", func() {
				Expect(config.port).To(Equal(8080))
			})
			It("host", func() {
				Expect(config.host).To(Equal("foobar"))
			})
			It("verbose", func() {
				Expect(config.verbose).To(Equal(true))
			})
			It("content", func() {
				Expect(config.content).To(Equal("boom"))
			})
			It("headers", func() {
				Expect(config.headers).To(Equal(headers))
			})
			It("deadTime", func() {
				Expect(config.deadTime).To(Equal(10 * time.Second))
			})
			It("minWait", func() {
				Expect(config.minWait).To(Equal(11 * time.Second))
			})
			It("maxWait", func() {
				Expect(config.maxWait).To(Equal(12 * time.Second))
			})
			It("randomWait", func() {
				Expect(config.randomWait).To(Equal(true))
			})
			It("minSize", func() {
				Expect(config.minSize).To(Equal(uint64(1000)))
			})
			It("maxSize", func() {
				Expect(config.maxSize).To(Equal(uint64(1000 * 2)))
			})
			It("randomSize", func() {
				Expect(config.randomSize).To(Equal(true))
			})
		})

	})

	Describe("ArgsConfigurationReader", func() {
		Describe("reads", func() {
			var args CommandLineArgs
			var config Configuration
			var headers []string
			BeforeEach(func() {
				headers = []string{"Age:1", "Content-type:text/plain"}

				args = CommandLineArgs{}
				args.Port = 8080
				args.Host = "foobar"
				args.Verbose = true
				args.Content = "boom"
				args.Headers = headers
				args.DeadTime = "10s"
				args.MinWait = "11s"
				args.MaxWait = "12s"
				args.RandomWait = true
				args.MinSize = "1KB"
				args.MaxSize = "2KB"
				args.RandomSize = true

				config = NewArgsConfigurationReader(&args).Read()
			})
			It("port", func() {
				Expect(config.port).To(Equal(8080))
			})
			It("host", func() {
				Expect(config.host).To(Equal("foobar"))
			})
			It("verbose", func() {
				Expect(config.verbose).To(Equal(true))
			})
			It("content", func() {
				Expect(config.content).To(Equal("boom"))
			})
			It("headers", func() {
				Expect(config.headers).To(Equal(headers))
			})
			It("deadTime", func() {
				Expect(config.deadTime).To(Equal(10 * time.Second))
			})
			It("minWait", func() {
				Expect(config.minWait).To(Equal(11 * time.Second))
			})
			It("maxWait", func() {
				Expect(config.maxWait).To(Equal(12 * time.Second))
			})
			It("randomWait", func() {
				Expect(config.randomWait).To(Equal(true))
			})
			It("minSize", func() {
				Expect(config.minSize).To(Equal(uint64(1000)))
			})
			It("maxSize", func() {
				Expect(config.maxSize).To(Equal(uint64(1000 * 2)))
			})
			It("randomSize", func() {
				Expect(config.randomSize).To(Equal(true))
			})
		})
	})
})
