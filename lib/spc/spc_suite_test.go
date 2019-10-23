package spc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spc Suite")
}
