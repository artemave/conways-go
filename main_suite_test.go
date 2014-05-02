package main_test

import (
	. "github.com/artemave/conways-go"
	. "github.com/artemave/conways-go/dependencies/ginkgo"
	. "github.com/artemave/conways-go/dependencies/gomega"

	"testing"
)

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RegisterRoutes()
	RunSpecs(t, "Main Suite")
}
