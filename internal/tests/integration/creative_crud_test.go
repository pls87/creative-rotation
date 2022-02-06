//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CreativeCRUDSuite struct {
	suite.Suite
	url string
}

func (s *CreativeCRUDSuite) SetupTest() {
	s.url = os.Getenv("CR_API_URL")
}

func (s *CreativeCRUDSuite) TearDownTest() {
}

func (s *CreativeCRUDSuite) NoopTest() {
}

func TestCreativeCRUD(t *testing.T) {
	suite.Run(t, new(HttpAPISuite))
}
