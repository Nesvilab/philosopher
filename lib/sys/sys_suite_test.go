package sys_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSys(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sys Suite")
}
