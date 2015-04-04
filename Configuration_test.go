package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Configuration", func() {

	Describe("ArgsConfigurationReader", func() {
		Describe("reads", func() {
			var args CommandLineArgs
			var config Configuration
			var headers []string
			Context("CommandLineArgs", func() {
				headers = []string{"Age:1", "Content-type:text/plain"}

				args = CommandLineArgs{}
				args.port = 8080
				args.host = "foobar"
				args.verbose = true
				args.content = "boom"
				args.headers = headers
				args.deadTime = "10s"
				args.minWait = "11s"
				args.maxWait = "12s"
				args.randomWait = true
				args.minSize = "1KB"
				args.maxSize = "2KB"
				args.randomSize = true

				config = NewArgsConfigurationReader(args).Read()
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
