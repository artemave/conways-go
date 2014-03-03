package routes_test

import (
	. "github.com/artemave/conways-go/dependencies/ginkgo"
	. "github.com/artemave/conways-go/dependencies/gomega"
	. "github.com/artemave/conways-go/routes"

	"testing"
)

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RegisterRoutes()
	RunSpecs(t, "Routes Suite")
}
