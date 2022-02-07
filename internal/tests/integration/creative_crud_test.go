//go:build integration
// +build integration

package integration

import (
	"net/http"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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
	if s.baseURL == "" {
		s.baseURL = "http://127.0.0.1:8081"
	}
	s.creatives = helpers.NewCreativeHelper(s.baseURL)
}

func (s *CreativeCRUDSuite) TearDownTest() {
}

func (s *CreativeCRUDSuite) TestCreateCreative() {
	car := gofakeit.Car()
	creativeDesc := car.Brand + " " + car.Model
	creative := s.newCreative(creativeDesc)

	s.Greaterf(creative.ID, 0, "ID of created creative should be more than 0")
	s.Equal(creativeDesc, creative.Desc)

	creatives := s.allCreatives()
	found := false
	for _, c := range creatives {
		if c.ID == creative.ID && c.Desc == creative.Desc {
			found = true
			break
		}
	}

	s.Truef(found, "created creative %v couldn't be found in storage", creative)
}

func (s *CreativeCRUDSuite) newCreative(desc string) (creative helpers.Entity) {
	code, resp, err := s.creatives.New(desc)
	s.NoErrorf(err, "no error expected but was: %s", err)
	s.Equal(http.StatusOK, code)
	creative, err = helpers.ParseOne("creative", resp)
	s.NoErrorf(err, "no error expected but was: %s", err)

	return creative
}

func (s *CreativeCRUDSuite) allCreatives() (creatives []helpers.Entity) {
	code, resp, err := s.creatives.All()
	s.NoErrorf(err, "no error expected but was: %s", err)
	s.Equal(http.StatusOK, code)
	creatives, err = helpers.ParseMany("creatives", resp)
	s.NoErrorf(err, "no error expected but was: %s", err)

	return creatives
}

func TestCreativeCRUD(t *testing.T) {
	suite.Run(t, new(CreativeCRUDSuite))
}
