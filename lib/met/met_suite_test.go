package met_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Met Suite")
}
