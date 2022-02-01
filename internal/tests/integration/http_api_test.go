//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HttpAPISuite struct {
	suite.Suite
	url string
}

func (s *HttpAPISuite) SetupTest() {
	s.url = os.Getenv("CR_API_URL")
	if s.url == "" {
		s.url = "http://127.0.0.1:8080"
	}
}

func (s *HttpAPISuite) NoopTest() {
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(HttpAPISuite))
}
