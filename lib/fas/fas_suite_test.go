package fas_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFas(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fas Suite")
}
