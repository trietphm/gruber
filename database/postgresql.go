package database

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/trietphm/gruber/config"
	"github.com/trietphm/gruber/model/mpg"
)

// PgI is a interface for manipulating data in database postgresql
type PgI interface {
	// CreateDriver Insert driver to database
	CreateDriver(driver *mpg.Driver) error

	// CreatePassenger Insert passenger to database
	CreatePassenger(passenger *mpg.Passenger) error

	// UpdateDriverState Update driver state in database
	UpdateDriverState(driverID int, state string) error

	// GetDriver get driver by id
	GetDriver(driverID int) (*mpg.Driver, error)
}

// Pg
type Pg struct {
	pg.DB
}

// OpenPostgresqlDB Open connection to postgresql
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
func (db *Pg) CreateDriver(driver *mpg.Driver) error {
	return db.Insert(driver)
}

// CreatePassenger Insert passenger to database
func (db *Pg) CreatePassenger(passenger *mpg.Passenger) error {
	return db.Insert(passenger)
}

// UpdateDriverState Update driver state in database
func (db *Pg) UpdateDriverState(driverID int, state string) error {
	_, err := db.Exec(`UPDATE drivers SET state = ? WHERE id = ?`, state, driverID)
	return err
}

// GetDriver Get driver by id
func (db *Pg) GetDriver(driverID int) (*mpg.Driver, error) {
	var driver mpg.Driver
	err := db.Model(&driver).Where("id = ?", driverID).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}

	return &driver, err
}
