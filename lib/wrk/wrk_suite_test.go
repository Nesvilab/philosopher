package wrk_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWrk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wrk Suite")
}
