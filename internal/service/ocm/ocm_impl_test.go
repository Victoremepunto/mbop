package ocm

import (
	"os"
	"testing"

	"github.com/redhatinsights/mbop/internal/config"

	v1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	is_internal_label string
}

func (suite *TestSuite) SetupSuite() {
	suite.is_internal_label = "foo"
}

func (suite *TestSuite) SetupTest() {
	config.Reset()
}

func (suite *TestSuite) TestGetIsInternalMatch() {
	os.Setenv("IS_INTERNAL_LABEL", "foo")
	l := &v1.LabelBuilder{}
	l.Value(suite.is_internal_label)
	acctB := &v1.AccountBuilder{}
	acctB.Labels(l)
	acct, _ := acctB.Build()
	assert.Equal(suite.T(), true, getIsInternal(acct))
}

func (suite *TestSuite) TestGetIsInternaEmptyLabels() {
	os.Setenv("IS_INTERNAL_LABEL", "foo")
	acctB := &v1.AccountBuilder{}
	acct, _ := acctB.Build()
	assert.Equal(suite.T(), false, getIsInternal(acct))
}

func (suite *TestSuite) TestGetIsInternalNoMatch() {
	os.Setenv("IS_INTERNAL_LABEL", "bar")
	l := &v1.LabelBuilder{}
	l.Value(suite.is_internal_label)
	acctB := &v1.AccountBuilder{}
	acctB.Labels(l)
	acct, _ := acctB.Build()
	assert.Equal(suite.T(), false, getIsInternal(acct))
}

func (suite *TestSuite) TearDownSuite() {
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
