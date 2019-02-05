package uti_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUti(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Uti Suite")
}
