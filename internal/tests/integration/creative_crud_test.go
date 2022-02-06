//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/pls87/creative-rotation/internal/tests/integration/helpers"
	"github.com/stretchr/testify/suite"
)

type CreativeCRUDSuite struct {
	suite.Suite
	creatives *helpers.CreativeHelper
	baseURL   string
}

func (s *CreativeCRUDSuite) SetupTest() {
	s.baseURL = os.Getenv("CR_API_URL")
	s.creatives = helpers.NewCreativeHelper(s.baseURL)
}

func (s *CreativeCRUDSuite) TearDownTest() {
}

func (s *CreativeCRUDSuite) TestCreateCreative() {
	code, resp, err := s.creatives.New("Jeep Wrangler")
	s.NoErrorf(err, "no error expected but was: %s", err)
	s.Equal(http.StatusOK, code)
	fmt.Println(string(resp))
}

func TestCreativeCRUD(t *testing.T) {
	suite.Run(t, new(CreativeCRUDSuite))
}
