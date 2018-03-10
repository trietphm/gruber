package model

import "time"

const (
	StateAvailable = "available"
	StateBusy      = "busy"
)

// Driver
type Driver struct {
	tableName struct{} `sql:"drivers,alias:drivers" pg:",discard_unknown_columns"`
	ID        int
	Name      string
	State     string
	CreatedAt time.Time
}

// Passenger
type Passenger struct {
	tableName struct{} `sql:"passengers,alias:passengers" pg:",discard_unknown_columns"`
	ID        int
	Name      string
	CreatedAt time.Time
}

// DriverLocation Location of driver which will be updated each 6s
type DriverLocation struct {
	tableName struct{} `sql:"driver_locations,alias:driver_locations" pg:",discard_unknown_columns"`
	ID        int
	DriverID  int
	Lat       float32
	Lng       float32
	CreatedAt time.Time
}
