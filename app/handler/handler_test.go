package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/trietphm/gruber/model/mcass"
	"github.com/trietphm/gruber/model/mpg"
	"github.com/trietphm/gruber/model/mredis"
)

type mockDbPg struct{}
type mockDbRedis struct{}
type mockDbCass struct{}

func newMockEngine(t *testing.T) *gin.Engine {
	mockDbPg := mockDbPg{}
	mockDbRedis := mockDbRedis{}
	mockDbCass := mockDbCass{}
	engine, err := NewEngine(mockDbPg, mockDbCass, mockDbRedis)
	if err != nil {
		t.FailNow()
		return nil
	}

	return engine
}

func TestGetDriverHistory(t *testing.T) {
	tt := []struct {
		url        string
		StatusCode int
		RespData   string
	}{
		{"/drivers/1/history", http.StatusOK, "[]"},
		{"/drivers/0/history", http.StatusNotFound, ""},
		{"/drivers/-1/history", http.StatusInternalServerError, `{"message":"INTERNAL SERVER ERROR"}`},
		{"/drivers/abc/history", http.StatusNotFound, ""},
	}

	router := newMockEngine(t)
	ts := httptest.NewServer(router)
	for _, tc := range tt {
		client := ts.Client()
		resp, err := client.Get(ts.URL + tc.url)
		if err != nil {
			t.Log(ts.URL+"/"+tc.url, err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("read data from resp body fail")
		}

		assert.Equal(t, tc.StatusCode, resp.StatusCode)
		assert.Equal(t, tc.RespData, string(body))
	}
	ts.Close()
}

func TestCreateDriver(t *testing.T) {
	tt := []struct {
		url        string
		input      string
		StatusCode int
		RespData   string
	}{
		{"/drivers", `{"name":"driver"}`, http.StatusOK, `{"id":1}`},
		{"/drivers", `{"name":1}`, http.StatusBadRequest, `{"message":"Invalid format"}`},
		{"/drivers", `{"name":""}`, http.StatusBadRequest, `{"message":"Name can not be empty"}`},
		{"/drivers", `{"name":"INVALID"}`, http.StatusInternalServerError, `{"message":"INTERNAL SERVER ERROR"}`},
	}

	router := newMockEngine(t)
	ts := httptest.NewServer(router)
	for _, tc := range tt {
		client := ts.Client()
		resp, err := client.Post(ts.URL+tc.url, "application/json", bytes.NewBuffer([]byte(tc.input)))
		if err != nil {
			t.Log(ts.URL+"/"+tc.url, err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("read data from resp body fail")
		}

		assert.Equal(t, tc.StatusCode, resp.StatusCode)
		assert.Equal(t, tc.RespData, string(body))
	}
	ts.Close()
}

func TestCreatePassenger(t *testing.T) {
	tt := []struct {
		url        string
		input      string
		StatusCode int
		RespData   string
	}{
		{"/passengers", `{"name":"passenger"}`, http.StatusOK, `{"id":1}`},
		{"/passengers", `{"name":1}`, http.StatusBadRequest, `{"message":"Invalid format"}`},
		{"/passengers", `{"name":""}`, http.StatusBadRequest, `{"message":"Name can not be empty"}`},
		{"/passengers", `{"name":"INVALID"}`, http.StatusInternalServerError, `{"message":"INTERNAL SERVER ERROR"}`},
	}

	router := newMockEngine(t)
	ts := httptest.NewServer(router)
	for _, tc := range tt {
		client := ts.Client()
		resp, err := client.Post(ts.URL+tc.url, "application/json", bytes.NewBuffer([]byte(tc.input)))
		if err != nil {
			t.Log(ts.URL+"/"+tc.url, err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("read data from resp body fail")
		}

		assert.Equal(t, tc.StatusCode, resp.StatusCode)
		assert.Equal(t, tc.RespData, string(body))
	}
	ts.Close()
}

func TestUpdateDriverLocation(t *testing.T) {
	tt := []struct {
		url        string
		input      string
		StatusCode int
		RespData   string
	}{
		{"/drivers/1/locations", `{"location":{"lat":30,"lng":100}}`, http.StatusOK, `{"id":1,"location":{"lat":30,"lng":100}}`},
		{"/drivers/abc/locations", `{"location":{"lat":30,"lng":100}}`, http.StatusNotFound, ``},
		{"/drivers/0/locations", `{"location":{"lat":30,"lng":100}}`, http.StatusNotFound, ``},
		{"/drivers/-1/locations", `{"location":{"lat":30,"lng":100}}`, http.StatusInternalServerError, `{"message":"INTERNAL SERVER ERROR"}`},
		{"/drivers/1/locations", `{"location":{"lat":"30a","lng":1aa00}}`, http.StatusBadRequest, `{"message":"Invalid format"}`},
	}

	router := newMockEngine(t)
	ts := httptest.NewServer(router)
	for _, tc := range tt {
		client := ts.Client()
		url := ts.URL + tc.url
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(tc.input)))
		if err != nil {
			t.Log(url, err)
			return
		}
		req.Header.Add("content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Log(url, err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("read data from resp body fail")
		}

		assert.Equal(t, tc.StatusCode, resp.StatusCode)
		assert.Equal(t, tc.RespData, string(body))
	}
	ts.Close()
}

func TestRequestDrivers(t *testing.T) {
	tt := []struct {
		url        string
		input      string
		StatusCode int
		RespData   string
	}{
		{"/requests", `{"passenger_id":1, "location":{"lat":30,"lng":100}}`, http.StatusOK, `[]`},
		{"/requests", `{"passenger_id":-1, "location":{"lat":30,"lng":100}}`, http.StatusBadRequest, `{"message":"Not found passenger"}`},
		{"/requests", `{"passenger_id":"1", "location":{"lat":30,"lng":100}}`, http.StatusBadRequest, `{"message":"Invalid format"}`},
	}

	router := newMockEngine(t)
	ts := httptest.NewServer(router)
	for _, tc := range tt {
		client := ts.Client()
		resp, err := client.Post(ts.URL+tc.url, "application/json", bytes.NewBuffer([]byte(tc.input)))
		if err != nil {
			t.Log(ts.URL+"/"+tc.url, err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("read data from resp body fail")
		}

		assert.Equal(t, tc.StatusCode, resp.StatusCode)
		assert.Equal(t, tc.RespData, string(body))
	}
	ts.Close()
}

func TestUpdateDriverState(t *testing.T) {
	tt := []struct {
		url        string
		input      string
		StatusCode int
		RespData   string
	}{
		{"/drivers/1", `{"state":"available"}`, http.StatusOK, `{}`},
		{"/drivers/1", `{"state":"busy"}`, http.StatusOK, `{}`},
		{"/drivers/1", `{"state":1}`, http.StatusBadRequest, `{"message":"Invalid format"}`},
		{"/drivers/1", `{"state":""}`, http.StatusBadRequest, `{"message":"Invalid state"}`},
		{"/drivers/1", `{"state":"abcd"}`, http.StatusBadRequest, `{"message":"Invalid state"}`},
		{"/drivers/0", `{"state":"busy"}`, http.StatusNotFound, ``},
		{"/drivers/-1", `{"state":"busy"}`, http.StatusInternalServerError, `{"message":"INTERNAL SERVER ERROR"}`},
	}

	router := newMockEngine(t)
	ts := httptest.NewServer(router)
	for _, tc := range tt {
		client := ts.Client()
		url := ts.URL + tc.url
		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer([]byte(tc.input)))
		if err != nil {
			t.Log(url, err)
			return
		}
		req.Header.Add("content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Log(url, err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("read data from resp body fail")
		}

		assert.Equal(t, tc.StatusCode, resp.StatusCode)
		assert.Equal(t, tc.RespData, string(body))
	}
	ts.Close()
}

// ===============================
// MOCK DATABASE
// ===============================

// CreateDriver Insert driver to database
func (mockDbPg) CreateDriver(driver *mpg.Driver) error {
	switch driver.Name {
	case "INVALID":
		return errors.New("Mock error database")
	default:
		driver.ID = 1
	}
	return nil

}

// CreatePassenger Insert passenger to database
func (mockDbPg) CreatePassenger(passenger *mpg.Passenger) error {
	switch passenger.Name {
	case "INVALID":
		return errors.New("Mock error database")
	default:
		passenger.ID = 1
	}
	return nil
}

// UpdateDriverState Update driver state in database
func (mockDbPg) UpdateDriverState(driverID int, state string) error {
	return nil
}

func (mockDbPg) GetDriver(driverID int) (*mpg.Driver, error) {
	switch driverID {
	case 1:
		return &mpg.Driver{
			ID:   1,
			Name: "Driver 1",
		}, nil
	case 0:
		return nil, nil
	case -1:
		return nil, errors.New("Mock db error")
	default:
		return &mpg.Driver{
			ID:   1,
			Name: "Driver 1",
		}, nil
	}
}

// =======
// Mock database redis
// =======

// UpdateLocation Update driver's location in database
func (mockDbRedis) PushDriverLocationGeo(driverID int, lat, lng float64) error {
	return nil
}

// Remove driver location from redis geo
func (mockDbRedis) RemoveDriverLocationGeo(driverID int) error {
	return nil
}

// GetNearestDrivers Get near available driver near a geo location
func (mockDbRedis) GetNearestDrivers(lat, lng, radius float64, limit int) ([]mredis.DriverLocation, error) {
	return []mredis.DriverLocation{}, nil
}

// =======
// Mock database cassandra
// =======

// DriverHistory Get driver history location

// UpdateLocation Update driver's location in database
func (mockDbCass) CreateDriverLocation(location *mcass.DriverLocation) error {
	location.DriverID = 1
	location.Lat = 30
	location.Lng = 100
	return nil
}

// DriverHistory Get driver history location
func (mockDbCass) GetDriverHistory(driverID int, from time.Time) ([]mcass.DriverLocation, error) {
	return []mcass.DriverLocation{}, nil
}

// GetDriverLatestLocation get latest driver location
func (mockDbCass) GetDriverLatestLocation(driverID int) (*mcass.DriverLocation, error) {
	switch driverID {
	case 1:
		return &mcass.DriverLocation{
			DriverID: 1,
			Lat:      45,
			Lng:      100,
		}, nil
	case 0:
		return nil, nil
	case -1:
		return nil, errors.New("Mock database error")
	}

	return nil, nil
}
