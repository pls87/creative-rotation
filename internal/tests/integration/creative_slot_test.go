//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"

	"github.com/pls87/creative-rotation/internal/tests/integration/helpers"
	"github.com/stretchr/testify/suite"
)

type CreativeSlotSuite struct {
	suite.Suite
	entities *helpers.EntityHelper
	baseURL  string
}

func (s *CreativeSlotSuite) SetupTest() {
	s.baseURL = os.Getenv("CR_API_URL")
	if s.baseURL == "" {
		s.baseURL = "http://127.0.0.1:8081"
	}
	s.entities = helpers.NewEntityHelper(s.baseURL)
}

func (s *CreativeSlotSuite) TearDownTest() {
}

func (s *CreativeSlotSuite) TestAddToSlot() {
}

func TestCreativeSlot(t *testing.T) {
	suite.Run(t, new(CreativeSlotSuite))
}
