package mpg

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
