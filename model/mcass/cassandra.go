package mcass

import "github.com/gocql/gocql"

// DriverLocation Location of driver which will be updated each 6s
type DriverLocation struct {
	ID        gocql.UUID
	DriverID  int
	Lat       float64
	Lng       float64
	CreatedAt int64
}
