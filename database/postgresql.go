package database

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/trietphm/gruber/config"
	"github.com/trietphm/gruber/model"
)

type Pg struct {
	pg.DB
}

func OpenPostgresqlDB(cfg config.Postgresql) (*Pg, error) {
	options := pg.Options{
		User:     cfg.User,
		Password: cfg.Password,
		Addr:     cfg.Host + ":" + cfg.Port,
		Database: cfg.Name,
	}
	db := pg.Connect(&options)

	// Try with simple query
	_, err := db.Exec("SELECT 1")
	if err != nil {
		panic(errors.Errorf("error on connecting to database %s", err))
	}

	return &Pg{*db}, err
}

// CreateDriver Insert driver to database
func (db *Pg) CreateDriver(driver *model.Driver) error {
	return db.Insert(driver)
}

// CreatePassenger Insert passenger to database
func (db *Pg) CreatePassenger(passenger *model.Passenger) error {
	return db.Insert(passenger)
}

// CreateDriverLocation Create new driver location in database
func (db *Pg) CreateDriverLocation(location *model.DriverLocation) error {
	return db.Insert(location)
}

// UpdateDriverState Update driver state in database
func (db *Pg) UpdateDriverState(driverID int, state string) error {
	_, err := db.Exec(`UPDATE drivers SET state = ? WHERE id = ?`, state, driverID)
	return err
}

// DriverHistory Get driver history location
func (db *Pg) GetDriverHistory(driverID int, from time.Time) ([]model.DriverLocation, error) {
	var locations []model.DriverLocation
	err := db.Model(&locations).
		Where("driver_id = ? AND created_at > ?", driverID, from).
		Select(&locations)
	return locations, err
}

// GetNearestDrivers Get near available driver near a geo location
func (db *Pg) GetNearestDrivers(lat, lng float32) ([]model.DriverLocation, error) {
	return []model.DriverLocation{}, nil
}

// GetDriver Get driver by id
func (db *Pg) GetDriver(driverID int) (*model.Driver, error) {
	var driver model.Driver
	err := db.Model(&driver).Where("id = ?", driverID).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}

	return &driver, err
}
