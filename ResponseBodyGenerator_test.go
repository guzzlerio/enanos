package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResponseBodyGenerator", func() {

	Describe("Max Response Body Generator", func() {
		It("generates a string of the defined lenth", func() {
			maxLength := 5
			generator := NewMaxResponseBodyGenerator(maxLength)
			value := generator.Generate()
			Expect(len(value)).To(Equal(maxLength))
		})
	})

	Describe("Random Response Body Generator", func() {
		It("generates a string of length between the defined min length and the defined max length", func() {
			minLength := 50
			maxLength := 500
			generator := NewRandomResponseBodyGenerator(minLength, maxLength)
			value := generator.Generate()
			Expect(len(value) >= minLength && len(value) <= maxLength).To(BeTrue())
		})
	})
})
