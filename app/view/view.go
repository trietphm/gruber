package view

import (
	"encoding/json"
	"time"

	"github.com/trietphm/gruber/model/mcass"
	"github.com/trietphm/gruber/model/mredis"
)

type timestamp time.Time

func (t timestamp) MarshalJSON() ([]byte, error) {
	str := time.Time(t).UTC().Format("2006-01-02T15:04:05Z")
	return json.Marshal(str)
}

func (t *timestamp) String() string {
	return time.Time(*t).String()
}

func (t *timestamp) Unix() int64 {
	return time.Time(*t).Unix()
}

// User response data when sign up user
type User struct {
	ID int `json:"id"`
}

// DriverRequest Response data when passenger request a ride
type DriverRequest struct {
	DriverID int      `json:"driver_id"`
	Location Location `json:"location"`
}

// DriverLocation Response updating driver location
type DriverLocation struct {
	ID       int      `json:"id"`
	Location Location `json:"location"`
}

// DriverHistory Response get driver history
type DriverHistory struct {
	Timestamp timestamp `json:"ts"`
	Location  Location  `json:"location"`
}

// Location Response geolocation with latitude, longitude
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// PopulateDriverRequests Populate repsonse for an array of drivers for ride request
func PopulateDriverRequests(drivers []mredis.DriverLocation) []DriverLocation {
	res := make([]DriverLocation, len(drivers))

	for i, driver := range drivers {
		res[i] = DriverLocation{
			ID: driver.DriverID,
			Location: Location{
				Lat: driver.Lat,
				Lng: driver.Lng,
			},
		}
	}

	return res
}

func PopulateDriverHistory(driverLocations []mcass.DriverLocation) []DriverHistory {
	resp := make([]DriverHistory, len(driverLocations))
	for i, driverLocation := range driverLocations {
		resp[i] = DriverHistory{
			Timestamp: timestamp(time.Unix(driverLocation.CreatedAt, 0)),
			Location: Location{
				Lat: driverLocation.Lat,
				Lng: driverLocation.Lng,
			},
		}
	}

	return resp
}
