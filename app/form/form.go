package form

import (
	"errors"

	"github.com/trietphm/gruber/model"
)

// Driver input for sign up driver
type Driver struct {
	Name string `json:"name"`
}

// Passenger input for sign up passenger
type Passenger struct {
	Name string `json:"name"`
}

// Location input location with latitude and longitude
type Location struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

// DriverState Input for update driver state
type DriverState struct {
	State string
}

// RequestRide Input for request a ride
type RequestRide struct {
	PassengerID int      `json:"passenger_id"`
	Location    Location `json:"location"`
}

// Validate Validate input request ride
func (input *RequestRide) Validate() error {
	if input.PassengerID <= 0 {
		return errors.New("Not found passenger")
	}

	if input.Location.Lat > 90 || input.Location.Lat < -90 {
		return errors.New("Invalid latitude")
	}

	if input.Location.Lng > 180 || input.Location.Lng < -180 {
		return errors.New("Invalid latitude")
	}

	return nil
}

// Validate validate update driver state
func (input *DriverState) Validate() error {
	if input.State == "" || (input.State != model.StateAvailable && input.State != model.StateBusy) {
		return errors.New("Invalid state")
	}

	return nil
}

// Validate validate input sign up driver
func (input *Driver) Validate() error {
	if input.Name == "" {
		return errors.New("Name can not be empty")
	}

	return nil
}

// Validate validate input sign up passenger
func (input *Passenger) Validate() error {
	if input.Name == "" {
		return errors.New("Name can not be empty")
	}

	return nil
}
