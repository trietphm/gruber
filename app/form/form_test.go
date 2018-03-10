package form

import "testing"

func TestRequestRideValidate(t *testing.T) {
	tt := []struct {
		input              RequestRide
		expectedErrMessage string
	}{
		{
			RequestRide{
				PassengerID: 1,
				Location: Location{
					Lat: 10.823099,
					Lng: 106.629664,
				},
			},
			"",
		},
		{
			RequestRide{
				PassengerID: 1,
				Location: Location{
					Lat: 100,
					Lng: 106.629664,
				},
			},
			"Invalid latitude",
		},
		{
			RequestRide{
				PassengerID: 1,
				Location: Location{
					Lat: 10.23423,
					Lng: 181,
				},
			},
			"Invalid longitude",
		},
		{
			RequestRide{
				PassengerID: 0,
				Location: Location{
					Lat: 10.23423,
					Lng: 101,
				},
			},
			"Not found passenger",
		},
		{
			RequestRide{
				PassengerID: -10,
				Location: Location{
					Lat: 10.23423,
					Lng: 101,
				},
			},
			"Not found passenger",
		},
	}
	for _, tc := range tt {
		err := tc.input.Validate()
		if tc.expectedErrMessage == "" {
			if err != nil {
				t.Errorf("FAIL with input: %v expected empty but output %s", tc.input, err.Error())
			}
		} else {
			if err.Error() != tc.expectedErrMessage {
				t.Errorf("FAIL with input: %v expected %s but output %s", tc.input, tc.expectedErrMessage, err.Error())
			}

		}
	}
}

func TestDriverStateValidate(t *testing.T) {
	tt := []struct {
		input              DriverState
		expectedErrMessage string
	}{
		{
			DriverState{
				State: "",
			},
			"Invalid state",
		},
		{
			DriverState{
				State: "availableeee",
			},
			"Invalid state",
		},
		{
			DriverState{
				State: "available",
			},
			"",
		},
		{
			DriverState{
				State: "busy",
			},
			"",
		},
	}
	for _, tc := range tt {
		err := tc.input.Validate()
		if tc.expectedErrMessage == "" {
			if err != nil {
				t.Errorf("FAIL with input: %v expected empty but output %s", tc.input, err.Error())
			}
		} else {
			if err.Error() != tc.expectedErrMessage {
				t.Errorf("FAIL with input: %v expected %s but output %s", tc.input, tc.expectedErrMessage, err.Error())
			}

		}
	}
}

func TestDriverValidate(t *testing.T) {
	tt := []struct {
		input              Driver
		expectedErrMessage string
	}{
		{
			Driver{
				Name: "",
			},
			"Name can not be empty",
		},
		{
			Driver{
				Name: "driver",
			},
			"",
		},
	}
	for _, tc := range tt {
		err := tc.input.Validate()
		if tc.expectedErrMessage == "" {
			if err != nil {
				t.Errorf("FAIL with input: %v expected empty but output %s", tc.input, err.Error())
			}
		} else {
			if err.Error() != tc.expectedErrMessage {
				t.Errorf("FAIL with input: %v expected %s but output %s", tc.input, tc.expectedErrMessage, err.Error())
			}

		}
	}
}

func TestPassengerValidate(t *testing.T) {
	tt := []struct {
		input              Passenger
		expectedErrMessage string
	}{
		{
			Passenger{
				Name: "",
			},
			"Name can not be empty",
		},
		{
			Passenger{
				Name: "driver",
			},
			"",
		},
	}
	for _, tc := range tt {
		err := tc.input.Validate()
		if tc.expectedErrMessage == "" {
			if err != nil {
				t.Errorf("FAIL with input: %v expected empty but output %s", tc.input, err.Error())
			}
		} else {
			if err.Error() != tc.expectedErrMessage {
				t.Errorf("FAIL with input: %v expected %s but output %s", tc.input, tc.expectedErrMessage, err.Error())
			}

		}
	}
}
