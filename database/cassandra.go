package database

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/trietphm/gruber/config"
	"github.com/trietphm/gruber/model/mcass"
)

// CassandraI is a interface for manipulating data in database cassandra
type CassandraI interface {
	// UpdateLocation Update driver's location in database
	CreateDriverLocation(location *mcass.DriverLocation) error

	// DriverHistory Get driver history location
	GetDriverHistory(driverID int, from time.Time) ([]mcass.DriverLocation, error)

	// GetDriverLatestLocation get latest driver location
	GetDriverLatestLocation(driverID int) (*mcass.DriverLocation, error)
}

type Cassandra struct {
	*gocql.Session
}

var _ CassandraI = CassandraI(Cassandra{})

// OpenPostgresqlDB Open connection to postgresql
func OpenCassandraDB(cfg config.Cassandra) (*Cassandra, error) {
	cluster := gocql.NewCluster([]string{cfg.Cluster}...)
	cluster.ProtoVersion = 4
	cluster.Timeout = 6 * time.Second
	cluster.Keyspace = cfg.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.User,
		Password: cfg.Password,
	}
	var err error
	casSession, err := cluster.CreateSession()
	return &Cassandra{casSession}, err
}

// CreateDriverLocation Create new driver location in database
func (db Cassandra) CreateDriverLocation(location *mcass.DriverLocation) error {
	err := db.Query(
		`INSERT INTO driver_locations("id", "driver_id", "created_at", "lat", "lng")
		VALUES (?, ?, ?, ?, ?) `,
		gocql.TimeUUID(), location.DriverID, location.CreatedAt, location.Lat, location.Lng).Exec()
	return err
}

// GetDriverHistory Get driver history
func (db Cassandra) GetDriverHistory(driverID int, from time.Time) ([]mcass.DriverLocation, error) {
	var locations []mcass.DriverLocation
	iter := db.Query(
		`SELECT "id", "driver_id", "created_at", "lat", "lng" 
		FROM  driver_locations
		WHERE driver_id = ? AND created_at > ?;`,
		driverID, from.Unix()).Iter()
	defer iter.Close()
	results, err := iter.SliceMap()
	if err != nil {
		return locations, err
	}
	locations = make([]mcass.DriverLocation, len(results))
	for i, row := range results {
		locations[i].ID = row["id"].(gocql.UUID)
		locations[i].DriverID = row["driver_id"].(int)
		locations[i].CreatedAt = row["created_at"].(int64)
		locations[i].Lat = row["lat"].(float64)
		locations[i].Lng = row["lng"].(float64)
	}
	return locations, nil
}

// GetDriverLatestLocation Get latest driver location
func (db Cassandra) GetDriverLatestLocation(driverID int) (*mcass.DriverLocation, error) {
	var locations []mcass.DriverLocation
	iter := db.Query(
		`SELECT "id", "driver_id", "created_at", "lat", "lng" 
		FROM  driver_locations
		WHERE driver_id = ?
		LIMIT 1;`,
		driverID).Iter()
	defer iter.Close()
	results, err := iter.SliceMap()
	if err != nil {
		return nil, err
	}
	locations = make([]mcass.DriverLocation, len(results))
	for i, row := range results {
		locations[i].ID = row["id"].(gocql.UUID)
		locations[i].DriverID = row["driver_id"].(int)
		locations[i].CreatedAt = row["created_at"].(int64)
		locations[i].Lat = row["lat"].(float64)
		locations[i].Lng = row["lng"].(float64)
	}
	if len(locations) == 0 {
		return nil, nil
	}
	location := locations[0]
	return &location, nil
}
