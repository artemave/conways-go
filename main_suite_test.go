package main_test

import (
	. "github.com/artemave/conways-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RegisterRoutes()
	RunSpecs(t, "Main Suite")
}
