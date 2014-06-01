package synchronized_broadcaster_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSynchronized_broadcaster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Synchronized_broadcaster Suite")
}
