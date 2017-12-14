package bio_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bio Suite")
}
