package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/logger"
	"github.com/redhatinsights/mbop/internal/store"
	"github.com/stretchr/testify/suite"
)

type AuthV1TestSuite struct {
	suite.Suite
	rec   *httptest.ResponseRecorder
	store store.Store
}

func (suite *AuthV1TestSuite) SetupSuite() {
	_ = logger.Init()
	config.Reset()
	os.Setenv("STORE_BACKEND", "memory")
}

func (suite *AuthV1TestSuite) BeforeTest(_, _ string) {
	suite.rec = httptest.NewRecorder()
	suite.Nil(store.SetupStore())

	// creating a new store for every test and overriding the dep injection function
	suite.store = store.GetStore()
	store.GetStore = func() store.Store { return suite.store }
}

func (suite *AuthV1TestSuite) AfterTest(_, _ string) {
	suite.rec.Result().Body.Close()
}

func (suite *AuthV1TestSuite) TestV1AuthNotFound() {
	req := httptest.NewRequest(http.MethodGet, "http://foobar/v1/auth", nil)
	req.Header.Set(CertHeader, "/CN=1234")
	AuthV1Handler(suite.rec, req)

	//nolint:bodyclose
	suite.Equal(http.StatusNotFound, suite.rec.Result().StatusCode)
}

func (suite *AuthV1TestSuite) TestV1AuthSuccess() {
	_, err := suite.store.Create(&store.Registration{OrgID: "12345", UID: "1234"})
	suite.Nil(err)

	req := httptest.NewRequest(http.MethodGet, "http://foobar/v1/auth", nil)
	req.Header.Set(CertHeader, "/CN=1234")
	AuthV1Handler(suite.rec, req)

	//nolint:bodyclose
	suite.Equal(http.StatusOK, suite.rec.Result().StatusCode)

	body, err := io.ReadAll(suite.rec.Body)
	suite.Nil(err)

	var resp AuthV1Response
	err = json.Unmarshal(body, &resp)
	suite.Nil(err)

	suite.Equal("cert", resp.Mechanism)
	suite.Equal("12345", resp.User.OrgID)
	suite.Equal("12345", resp.User.DisplayName)
	suite.Equal(-1, resp.User.ID)
	suite.Equal(true, resp.User.IsActive)
	suite.Equal(true, resp.User.IsOrgAdmin)
	suite.Equal("system", resp.User.Type)
}
