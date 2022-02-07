//go:build integration
// +build integration

package integration

import (
	"fmt"
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

func (s *CreativeCRUDSuite) TestCreateEmptyCreative() {
	s.testCreateEmptyEntity("creative")
}

func (s *CreativeCRUDSuite) TestCreateEmptySlot() {
	s.testCreateEmptyEntity("slot")
}

func (s *CreativeCRUDSuite) TestCreateEmptySegment() {
	s.testCreateEmptyEntity("segment")
}

func (s *CreativeCRUDSuite) testCreateEntity(kind string) {
	desc := gofakeit.BuzzWord()
	entity := s.entities.New(s.T(), kind, desc)

	s.Greaterf(entity.ID, 0, fmt.Sprintf("ID of created %s should be more than 0", kind))
	s.Equal(desc, entity.Desc)

	entities := s.entities.All(s.T(), kind)
	found := false
	for _, e := range entities {
		if e.ID == entity.ID && e.Desc == entity.Desc {
			found = true
			break
		}
	}

	s.Truef(found, "created %s %v couldn't be found in storage", kind, entity)
}

func (s *CreativeCRUDSuite) testCreateEmptyEntity(kind string) {
	code, _, err := s.entities.Push(kind, "")

	s.NoErrorf(err, "no error expected but got", err)
	s.Equal(code, http.StatusBadRequest)

	entities := s.entities.All(s.T(), kind)
	found := false
	for _, e := range entities {
		if e.Desc == "" {
			found = true
			break
		}
	}

	s.Falsef(found, "%s with empty description found", kind)
}

func TestCreativeCRUD(t *testing.T) {
	suite.Run(t, new(CreativeCRUDSuite))
}
