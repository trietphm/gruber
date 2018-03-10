package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trietphm/gruber/model"
)

type MockDatabase struct{}

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

	mockDB := MockDatabase{}
	router, err := NewEngine(mockDB)
	if err != nil {
		t.FailNow()
		return
	}

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

	mockDB := MockDatabase{}
	router, err := NewEngine(mockDB)
	if err != nil {
		t.FailNow()
		return
	}

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

	mockDB := MockDatabase{}
	router, err := NewEngine(mockDB)
	if err != nil {
		t.FailNow()
		return
	}

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
		{"/drivers/1/locations", `{"location":{"lat":30,"lng":100}}`, http.StatusOK, `{"id":1,"location":{"lat":45,"lng":100}}`},
		{"/drivers/abc/locations", `{"location":{"lat":30,"lng":100}}`, http.StatusNotFound, ``},
		{"/drivers/0/locations", `{"location":{"lat":30,"lng":100}}`, http.StatusNotFound, ``},
		{"/drivers/-1/locations", `{"location":{"lat":30,"lng":100}}`, http.StatusInternalServerError, `{"message":"INTERNAL SERVER ERROR"}`},
		{"/drivers/1/locations", `{"location":{"lat":"30a","lng":1aa00}}`, http.StatusBadRequest, `{"message":"Invalid format"}`},
	}

	mockDB := MockDatabase{}
	router, err := NewEngine(mockDB)
	if err != nil {
		t.FailNow()
		return
	}

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

	mockDB := MockDatabase{}
	router, err := NewEngine(mockDB)
	if err != nil {
		t.FailNow()
		return
	}

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

	mockDB := MockDatabase{}
	router, err := NewEngine(mockDB)
	if err != nil {
		t.FailNow()
		return
	}

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

// CreateDriver Insert driver to database
func (MockDatabase) CreateDriver(driver *model.Driver) error {
	switch driver.Name {
	case "INVALID":
		return errors.New("Mock error database")
	default:
		driver.ID = 1
	}
	return nil

}

// CreatePassenger Insert passenger to database
func (MockDatabase) CreatePassenger(passenger *model.Passenger) error {
	switch passenger.Name {
	case "INVALID":
		return errors.New("Mock error database")
	default:
		passenger.ID = 1
	}
	return nil
}

// UpdateDriverState Update driver state in database
func (MockDatabase) UpdateDriverState(driverID int, state string) error {
	return nil
}

func (MockDatabase) GetDriver(driverID int) (*model.Driver, error) {
	switch driverID {
	case 1:
		return &model.Driver{
			ID:   1,
			Name: "Driver 1",
		}, nil
	case 0:
		return nil, nil
	case -1:
		return nil, errors.New("Mock db error")
	default:
		return &model.Driver{
			ID:   1,
			Name: "Driver 1",
		}, nil
	}
}

// UpdateLocation Update driver's location in database
func (MockDatabase) CreateDriverLocation(location *model.DriverLocation) error {
	location.ID = 1
	location.Lat = 45
	location.Lng = 100
	return nil
}

// DriverHistory Get driver history location
func (MockDatabase) GetDriverHistory(driverID int, from time.Time) ([]model.DriverLocation, error) {
	return []model.DriverLocation{}, nil
}

// GetNearestDrivers Get near available driver near a geo location
func (MockDatabase) GetNearestDrivers(lat, lng float32) ([]model.DriverLocation, error) {
	return []model.DriverLocation{}, nil
}
