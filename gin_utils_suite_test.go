package ginUtils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGinUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GinUtils Suite")
}
