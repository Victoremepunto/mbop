package catchall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/logger"
	"github.com/redhatinsights/mbop/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegistrationTestSuite struct {
	suite.Suite
	rec   *httptest.ResponseRecorder
	store store.Store
}

func (suite *RegistrationTestSuite) SetupSuite() {
	_ = logger.Init()
	config.Reset()
	os.Setenv("STORE_BACKEND", "memory")
}

func (suite *RegistrationTestSuite) BeforeTest(_, _ string) {
	suite.rec = httptest.NewRecorder()
	suite.Nil(store.SetupStore())

	// creating a new store for every test and overriding the dep injection function
	suite.store = store.GetStore()
	store.GetStore = func() store.Store { return suite.store }
}

func (suite *RegistrationTestSuite) AfterTest(_, _ string) {
	suite.rec.Result().Body.Close()
}

func (suite *RegistrationTestSuite) TestJWTGet() {
	testData, _ := os.ReadFile("testdata/jwt.json")
	testDataStruct := &JSONStruct{}
	err := json.Unmarshal([]byte(testData), testDataStruct)
	assert.Nil(suite.T(), err, "error was not nil")

	k8sServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/realms/redhat-external/" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(testData)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer k8sServer.Close()

	os.Setenv("KEYCLOAK_SERVER", k8sServer.URL)

	// dummy muxer for the test
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(MakeNewMBOPServer().MainHandler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	resp, err := http.Get(fmt.Sprintf("%s/v1/jwt", sut.URL))
	b, _ := io.ReadAll(resp.Body)

	assert.Nil(suite.T(), err, "error was not nil")
	assert.Equal(suite.T(), 200, resp.StatusCode, "status code not good")
	assert.Equal(suite.T(), testDataStruct.PublicKey, string(b), fmt.Sprintf("expected body doesn't match: %v", string(b)))

	defer resp.Body.Close()
}

func (suite *RegistrationTestSuite) TestGetUrl() {
	os.Setenv("KEYCLOAK_SERVER", "http://test")
	path := MakeNewMBOPServer().getURL("path", map[string]string{"hi": "you"})
	assert.Equal(suite.T(), "http://test/path?hi=you", path, "did not match")
}

func (suite *RegistrationTestSuite) TearDownSuite() {
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrationTestSuite))
}

func (suite *RegistrationTestSuite) TestGoodRegistration() {
	store.SetupStore()
	db := store.GetStore()
	db.Create(&store.Registration{
		ID:          "nark",
		OrgID:       "12345",
		Username:    "nark",
		UID:         "nark",
		DisplayName: "foobar",
		Extra:       map[string]interface{}{},
		CreatedAt:   time.Time{},
	})

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(MakeNewMBOPServer().MainHandler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	req2, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/check_registration", sut.URL), nil)
	req2.Header.Set("x-rh-check-reg", "nark")
	resp, err := http.DefaultClient.Do(req2)
	suite.Nil(err)

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *RegistrationTestSuite) TestBadRegistration() {

	store.SetupStore()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(MakeNewMBOPServer().MainHandler))

	sut := httptest.NewServer(mux)
	defer sut.Close()

	req2, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/check_registration", sut.URL), nil)
	req2.Header.Set("x-rh-check-reg", "nark")
	resp, err := http.DefaultClient.Do(req2)
	suite.Nil(err)

	suite.Equal(http.StatusForbidden, resp.StatusCode)
}
