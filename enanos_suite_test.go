package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEnanos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enanos Suite")
}
