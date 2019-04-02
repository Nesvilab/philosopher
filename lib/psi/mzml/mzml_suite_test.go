package mzml_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMzn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MzML Suite")
}
