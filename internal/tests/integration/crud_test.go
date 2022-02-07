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
	entities *helpers.EntityHelper
	baseURL  string
}

func (s *CreativeCRUDSuite) SetupTest() {
	s.baseURL = os.Getenv("CR_API_URL")
	if s.baseURL == "" {
		s.baseURL = "http://127.0.0.1:8081"
	}
	s.entities = helpers.NewEntityHelper(s.baseURL)
}

func (s *CreativeCRUDSuite) TearDownTest() {
}

func (s *CreativeCRUDSuite) TestCreateCreative() {
	s.testCreateEntity("creative")
}

func (s *CreativeCRUDSuite) TestCreateSlot() {
	s.testCreateEntity("slot")
}

func (s *CreativeCRUDSuite) TestCreateSegment() {
	s.testCreateEntity("segment")
}

func (s *CreativeCRUDSuite) testCreateEntity(t string) {
	desc := gofakeit.BuzzWord()
	entity := s.new(t, desc)

	s.Greaterf(entity.ID, 0, "ID of created entity should be more than 0")
	s.Equal(desc, entity.Desc)

	entities := s.all(t)
	found := false
	for _, e := range entities {
		if e.ID == entity.ID && e.Desc == entity.Desc {
			found = true
			break
		}
	}

	s.Truef(found, "created entity %v couldn't be found in storage", entity)
}

func (s *CreativeCRUDSuite) new(t, desc string) (entity helpers.Entity) {
	code, resp, err := s.entities.New(t, desc)
	s.NoErrorf(err, "no error expected but was: %s", err)
	s.Equal(http.StatusOK, code)
	entity, err = helpers.ParseOne(t, resp)
	s.NoErrorf(err, "no error expected but was: %s", err)

	return entity
}

func (s *CreativeCRUDSuite) all(t string) (entities []helpers.Entity) {
	code, resp, err := s.entities.All(t)
	s.NoErrorf(err, "no error expected but was: %s", err)
	s.Equal(http.StatusOK, code)
	entities, err = helpers.ParseMany(t+"s", resp)
	s.NoErrorf(err, "no error expected but was: %s", err)

	return entities
}

func TestCreativeCRUD(t *testing.T) {
	suite.Run(t, new(CreativeCRUDSuite))
}
