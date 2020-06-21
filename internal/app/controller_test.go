package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"geogame/internal/locations"
	"geogame/internal/middleware"
	"geogame/internal/players"
)

type testControllerSuite struct {
	suite.Suite
	recorder *httptest.ResponseRecorder
	router   chi.Router
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(testControllerSuite))
}

func (suite *testControllerSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.router = chi.NewRouter()
	locStore := locations.NewMockStore()
	locationsSvc := locations.NewDefaultService(zap.NewNop(), locStore)

	playersStore := players.NewMockStore()
	playersSvc := players.NewDefaultService(zap.NewNop(), playersStore, time.Second*3, "")

	mockAuther := middleware.NewMockAuther()
	controller := NewController(zap.NewNop(), locationsSvc, playersSvc, mockAuther)
	controller.SetupRouter(suite.router)
}

func (suite *testControllerSuite) TestController_CreateLocation() {
	req := suite.Require()
	payload := locations.Location{
		ID: "1",
		GeoPoint: locations.GeoPoint{
			Longitude: 10.1,
			Latitude:  10.1,
		},
		MetaData: locations.MetaData{
			LocationName: "locationName",
			LocationType: "locationType",
		},
	}
	bs, err := json.Marshal(payload)
	req.NoError(err)

	buffer := bytes.NewBuffer(bs)
	request := httptest.NewRequest("POST", "/admin/loc/create", buffer)

	suite.router.ServeHTTP(suite.recorder, request)
	response := suite.recorder.Result()
	req.Equal(http.StatusOK, response.StatusCode)
}

func (suite *testControllerSuite) TestController_GetLocation() {
	req := suite.Require()

	request := httptest.NewRequest("GET", "/admin/loc/1", nil)

	suite.router.ServeHTTP(suite.recorder, request)
	response := suite.recorder.Result()
	req.Equal(http.StatusOK, response.StatusCode)
}

func (suite *testControllerSuite) TestController_UpdateLocation() {
	req := suite.Require()
	payload := locations.Location{
		ID: "1",
		GeoPoint: locations.GeoPoint{
			Longitude: 10.1,
			Latitude:  10.1,
		},
		MetaData: locations.MetaData{
			LocationName: "locationName",
			LocationType: "locationType",
		},
	}
	bs, err := json.Marshal(payload)
	req.NoError(err)

	buffer := bytes.NewBuffer(bs)
	request := httptest.NewRequest("PUT", "/admin/loc/update", buffer)

	suite.router.ServeHTTP(suite.recorder, request)
	response := suite.recorder.Result()
	req.Equal(http.StatusOK, response.StatusCode)
}

func (suite *testControllerSuite) TestController_DeleteLocation() {
	req := suite.Require()

	request := httptest.NewRequest("DELETE", "/admin/loc/1/delete", nil)

	suite.router.ServeHTTP(suite.recorder, request)
	response := suite.recorder.Result()
	req.Equal(http.StatusOK, response.StatusCode)
}

func (suite *testControllerSuite) TestController_RegisterClient() {
	req := suite.Require()
	payload := players.RegisterPayload{
		Name:     "dummyname",
		Email:    "dummy@mail.com",
		Password: "dummypassword",
	}
	bs, err := json.Marshal(payload)
	req.NoError(err)

	buffer := bytes.NewBuffer(bs)
	request := httptest.NewRequest("POST", "/client/register", buffer)

	suite.router.ServeHTTP(suite.recorder, request)
	response := suite.recorder.Result()
	req.Equal(http.StatusOK, response.StatusCode)
}

func (suite *testControllerSuite) TestController_LoginClient() {
	req := suite.Require()
	payload := players.LoginPayload{
		Email:    "dummy@mail.com",
		Password: "dummypassword",
	}
	bs, err := json.Marshal(payload)
	req.NoError(err)

	buffer := bytes.NewBuffer(bs)
	request := httptest.NewRequest("POST", "/client/login", buffer)

	suite.router.ServeHTTP(suite.recorder, request)
	response := suite.recorder.Result()
	req.Equal(http.StatusInternalServerError, response.StatusCode)
}

func (suite *testControllerSuite) TestController_SendLocation() {
	req := suite.Require()
	payload := locations.Location{
		ID: "1",
		GeoPoint: locations.GeoPoint{
			Longitude: 10.1,
			Latitude:  10.1,
		},
		MetaData: locations.MetaData{
			LocationName: "locationName",
			LocationType: "locationType",
		},
	}
	bs, err := json.Marshal(payload)
	req.NoError(err)

	buffer := bytes.NewBuffer(bs)
	request := httptest.NewRequest("POST", "/client/loc/Send", buffer)

	suite.router.ServeHTTP(suite.recorder, request)
	response := suite.recorder.Result()
	req.Equal(http.StatusNotFound, response.StatusCode)
}
