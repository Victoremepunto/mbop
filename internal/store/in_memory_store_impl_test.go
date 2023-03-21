package store

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type InMemoryStoreTestSuite struct {
	suite.Suite
	store Store
}

func (suite *InMemoryStoreTestSuite) SetupSuite() {}

func (suite *InMemoryStoreTestSuite) TearDownSuite() {}

func (suite *InMemoryStoreTestSuite) BeforeTest(_, _ string) {
	suite.store = &inMemoryStore{db: make([]Registration, 0)}
}

func TestSuiteRunInMemoryStore(t *testing.T) {
	suite.Run(t, new(InMemoryStoreTestSuite))
}

func (suite *InMemoryStoreTestSuite) TestCreate() {
	_, err := suite.store.Create(&Registration{})
	suite.Nil(err)
}

func (suite *InMemoryStoreTestSuite) TestFind() {
	_, err := suite.store.Create(&Registration{OrgID: "1234", UID: "1234"})
	suite.Nil(err)
	_, err = suite.store.Find("1234", "1234")
	suite.Nil(err)
}

func (suite *InMemoryStoreTestSuite) TestFindNotThere() {
	_, err := suite.store.Find("1234", "1234")
	suite.Error(err)
}

func (suite *InMemoryStoreTestSuite) TestFindByUID() {
	_, err := suite.store.Create(&Registration{OrgID: "1234", UID: "1234"})
	suite.Nil(err)
	_, err = suite.store.FindByUID("1234")
	suite.Nil(err)
}

func (suite *InMemoryStoreTestSuite) TestFindByUIDNotThere() {
	_, err := suite.store.FindByUID("1234")
	suite.Error(err)
}

func (suite *InMemoryStoreTestSuite) TestAll() {
	_, err := suite.store.Create(&Registration{OrgID: "1234", UID: "1234", DisplayName: "one"})
	suite.Nil(err)
	_, err = suite.store.Create(&Registration{OrgID: "2345", UID: "2345", DisplayName: "two"})
	suite.Nil(err)

	_, count, err := suite.store.All("1234", 0, 0)
	suite.Nil(err)

	suite.Equal(count, 1)
}

func (suite *InMemoryStoreTestSuite) TestDelete() {
	_, err := suite.store.Create(&Registration{OrgID: "1234"})
	suite.Nil(err)

	err = suite.store.Delete("1234", "")
	suite.Nil(err)
}

func (suite *InMemoryStoreTestSuite) TestDeleteNotExisting() {
	err := suite.store.Delete("1234", "")
	suite.Error(err)
}
