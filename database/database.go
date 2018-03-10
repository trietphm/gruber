package database

import (
	"time"

	"github.com/trietphm/gruber/model"
)

// Database is interface for manipulating data in database
type Database interface {
	// CreateDriver Insert driver to database
	CreateDriver(driver *model.Driver) error

	// CreatePassenger Insert passenger to database
	CreatePassenger(passenger *model.Passenger) error

	// UpdateDriverState Update driver state in database
	UpdateDriverState(driverID int, state string) error

	GetDriver(driverID int) (*model.Driver, error)

	// UpdateLocation Update driver's location in database
	CreateDriverLocation(location *model.DriverLocation) error

	// DriverHistory Get driver history location
	GetDriverHistory(driverID int, from time.Time) ([]model.DriverLocation, error)

	// GetNearestDrivers Get near available driver near a geo location
	GetNearestDrivers(lat, lng float32) ([]model.DriverLocation, error)
}
