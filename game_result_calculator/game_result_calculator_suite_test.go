package game_result_calculator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGameResultCalculator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Game Result Calculator Suite")
}
