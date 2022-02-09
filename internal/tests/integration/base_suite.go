//go:build integration
// +build integration

package integration

import (
	"os"

	"github.com/pls87/creative-rotation/internal/tests/integration/sql"
	"github.com/stretchr/testify/suite"
)

type BaseSuite struct {
	suite.Suite
	baseURL string
	client  sql.Client
}

func (s *BaseSuite) SetupSuite() {
	s.client = sql.Client{}
	s.baseURL = os.Getenv("CR_API_URL")
	if s.baseURL == "" {
		s.baseURL = "http://127.0.0.1:8081"
	}
	s.NoError(s.client.Connect())
	s.NoError(s.client.RunFile("./sql/clean.sql"))
}

func (s *BaseSuite) TearDownSuite() {
	s.NoError(s.client.RunFile("./sql/clean.sql"))
	s.NoError(s.client.Close())
}
