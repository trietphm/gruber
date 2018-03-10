package form

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
