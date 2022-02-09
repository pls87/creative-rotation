//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pls87/creative-rotation/internal/tests/integration/helpers"
	"github.com/stretchr/testify/suite"
)

type CRUDSuite struct {
	BaseSuite
	entities *helpers.EntityHelper
}

func (s *CRUDSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()
	s.entities = helpers.NewEntityHelper(s.baseURL)
}

func (s *CRUDSuite) TestCreateCreative() {
	s.testCreateEntity("creative")
}

func (s *CRUDSuite) TestCreateSlot() {
	s.testCreateEntity("slot")
}

func (s *CRUDSuite) TestCreateSegment() {
	s.testCreateEntity("segment")
}

func (s *CRUDSuite) TestCreateEmptyCreative() {
	s.testCreateEmptyEntity("creative")
}

func (s *CRUDSuite) TestCreateEmptySlot() {
	s.testCreateEmptyEntity("slot")
}

func (s *CRUDSuite) TestCreateEmptySegment() {
	s.testCreateEmptyEntity("segment")
}

func (s *CRUDSuite) testCreateEntity(kind string) {
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

func (s *CRUDSuite) testCreateEmptyEntity(kind string) {
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
	suite.Run(t, new(CRUDSuite))
}
