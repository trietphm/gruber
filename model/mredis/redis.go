package mredis

const (
	KeyDriverGeo = "DRIVER_GEO"
)

type DriverLocation struct {
	DriverID int
	Lat      float64
	Lng      float64
}
