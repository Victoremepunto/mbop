package mailer

import (
	"context"
	"os"
	"testing"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/logger"
	"github.com/redhatinsights/mbop/internal/models"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupSuite() {
	config.Reset()
	os.Setenv("MAILER_MODULE", "print")
	os.Setenv("USERS_MODULE", "mock")

	_ = logger.Init()
}

func (suite *TestSuite) TestMockConversionAll() {
	email := models.Email{
		Recipients: []string{"me"},
		CcList:     []string{"you"},
		BccList:    []string{"everyone"},
	}
	err := LookupEmailsForUsernames(context.Background(), &email)
	suite.Nil(err)
	suite.Equal("me@mocked.biz", email.Recipients[0])
	suite.Equal("you@mocked.biz", email.CcList[0])
	suite.Equal("everyone@mocked.biz", email.BccList[0])
}

func (suite *TestSuite) TestMockConversionSome() {
	email := models.Email{
		Recipients: []string{"me"},
		CcList:     []string{"you@gmail.com"},
	}
	err := LookupEmailsForUsernames(context.Background(), &email)
	suite.Nil(err)
	suite.Equal("me@mocked.biz", email.Recipients[0])
	suite.Equal("you@gmail.com", email.CcList[0])
}
