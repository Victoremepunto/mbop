package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/logger"

	"github.com/RedHatInsights/jwk2pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	testData       []byte
	testDataStruct *jwk2pem.JWKeys
	mockServer     *httptest.Server
	testPem        []byte
}

func (suite *TestSuite) SetupSuite() {
	_ = logger.Init()
	suite.testData, _ = os.ReadFile("testdata/jwks.json")
	suite.testDataStruct = &jwk2pem.JWKeys{}
	err := json.Unmarshal([]byte(suite.testData), suite.testDataStruct)
	assert.Nil(suite.T(), err, "error was not nil")
	suite.testPem, _ = os.ReadFile("testdata/pem.json")
}

func (suite *TestSuite) TestAwsJWTGetNoKid() {
	suite.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(suite.testData)
	}))
	defer suite.mockServer.Close()
	config.Reset()

	os.Setenv("JWT_MODULE", "aws")
	os.Setenv("JWK_URL", fmt.Sprintf("%s/v1/jwt", suite.mockServer.URL))

	// dummy muxer for the test
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(JWTV1Handler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/jwt", sut.URL))
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(suite.T(), err, "error was not nil")
	assert.Equal(suite.T(), 400, resp.StatusCode, "status code not good")
	assert.Equal(suite.T(), "{\"message\":\"kid required to return correct pub key\"}", string(b), fmt.Sprintf("expected body doesn't match: %v", string(b)))

	defer resp.Body.Close()
}

func (suite *TestSuite) TestAwsJWTGetNoKidMatch() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(suite.testData)
	}))
	defer mockServer.Close()
	config.Reset()

	os.Setenv("JWT_MODULE", "aws")
	os.Setenv("JWK_URL", fmt.Sprintf("%s/v1/jwt", mockServer.URL))
	kid := "123"

	// dummy muxer for the test
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(JWTV1Handler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/jwt?kid=%s", sut.URL, kid))
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(suite.T(), err, "error was not nil")
	assert.Equal(suite.T(), 404, resp.StatusCode, "status code not good")
	assert.Equal(suite.T(), "{\"message\":\"no JWK for kid: 123\"}", string(b), fmt.Sprintf("expected body doesn't match: %v", string(b)))

	defer resp.Body.Close()
}

func (suite *TestSuite) TestAwsJWTGetKidMatch() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(suite.testData)
	}))
	defer mockServer.Close()
	config.Reset()

	os.Setenv("JWT_MODULE", "aws")
	os.Setenv("JWK_URL", fmt.Sprintf("%s/v1/jwt", mockServer.URL))
	kid := "b4OUzJFABPSRwxX5VN7lYswVj9qoc3tet0tsfG5MSME"

	// dummy muxer for the test
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(JWTV1Handler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/jwt?kid=%s", sut.URL, kid))
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(suite.T(), err, "error was not nil")
	assert.Equal(suite.T(), 200, resp.StatusCode, "status code not good")
	assert.Equal(suite.T(), string(suite.testPem), string(b), fmt.Sprintf("expected body doesn't match: %v", string(b)))

	defer resp.Body.Close()
}

func (suite *TestSuite) TearDownSuite() {
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
