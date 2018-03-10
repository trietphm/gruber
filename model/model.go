package model

import "time"

const (
	StateAvailable = "available"
	StateBusy      = "busy"
)

// Driver
type Driver struct {
	ID        int
	Name      string
	State     string
	CreatedAt time.Time
}

// Passenger
type Passenger struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

// DriverLocation Location of driver which will be updated each 6s
type DriverLocation struct {
	ID        int
	DriverID  int
	Lat       float32
	Lng       float32
	CreatedAt time.Time
}
