package view

import (
	"time"

	"github.com/trietphm/gruber/model"
)

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
	Timestamp time.Time `json:"ts"`
	Location  Location  `json:"location"`
}

// Location Response geolocation with latitude, longitude
type Location struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

// PopulateDriverRequests Populate repsonse for an array of drivers for ride request
func PopulateDriverRequests(drivers []model.DriverLocation) []DriverLocation {
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

func PopulateDriverHistory(driverLocations []model.DriverLocation) []DriverHistory {
	resp := make([]DriverHistory, len(driverLocations))
	for i, driverLocation := range driverLocations {
		resp[i] = DriverHistory{
			Timestamp: driverLocation.CreatedAt,
			Location: Location{
				Lat: driverLocation.Lat,
				Lng: driverLocation.Lng,
			},
		}
	}

	return resp
}
