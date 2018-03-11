package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trietphm/gruber/app/form"
	"github.com/trietphm/gruber/app/view"
	"github.com/trietphm/gruber/database"
	"github.com/trietphm/gruber/model/mcass"
	"github.com/trietphm/gruber/model/mpg"
	"github.com/trietphm/gruber/util"
)

// Handler Manager handler functions
type Handler struct {
	dbPg    database.PgI
	dbCass  database.CassandraI
	dbRedis database.RedisI
}

// NewEngine Setup API router
func NewEngine(dbPg database.PgI, dbCass database.CassandraI, dbRedis database.RedisI) (*gin.Engine, error) {
	engine := gin.Default()
	handler := Handler{
		dbPg:    dbPg,
		dbCass:  dbCass,
		dbRedis: dbRedis,
	}
	router := engine.Group("")
	router.POST("/passengers", handler.CreatePassenger)
	router.POST("/requests", handler.RequestDrivers)

	driverGroup := router.Group("/drivers")
	driverGroup.POST("", handler.CreateDriver)
	driverGroup.PUT("/:id/locations", handler.UpdateDriverLocation)
	driverGroup.GET("/:id/history", handler.GetDriverHistory)
	driverGroup.PATCH("/:id", handler.UpdateDriverState)

	return engine, nil
}

// CreatePassenger Sign up passenger
func (h *Handler) CreatePassenger(c *gin.Context) {
	var input form.Passenger
	if err := c.Bind(&input); err != nil {
		util.RespBadRequest(c, "Invalid format")
		return
	}

	if err := input.Validate(); err != nil {
		util.RespBadRequest(c, err.Error())
		return
	}

	passenger := mpg.Passenger{
		Name: input.Name,
	}
	if err := h.dbPg.CreatePassenger(&passenger); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	resp := view.User{
		ID: passenger.ID,
	}
	util.RespOK(c, resp)
}

// RequestDrivers Get list nearest drivers
func (h *Handler) RequestDrivers(c *gin.Context) {
	var input form.RequestRide
	if err := c.Bind(&input); err != nil {
		util.RespBadRequest(c, "Invalid format")
		return
	}

	if err := input.Validate(); err != nil {
		util.RespBadRequest(c, err.Error())
		return
	}

	// Set default radius 5 km. TODO if not enough then increase radius
	var defaultRadius float64 = 5
	numberOfTop := 5
	drivers, err := h.dbRedis.GetNearestDrivers(input.Location.Lat, input.Location.Lng, defaultRadius, numberOfTop)
	if err != nil {
		util.RespInternalServerError(c, err)
		return
	}
	fmt.Println(drivers)

	resp := view.PopulateDriverRequests(drivers)
	util.RespOK(c, resp)
}

// CreateDriver Sign up driver
func (h *Handler) CreateDriver(c *gin.Context) {
	var input form.Driver
	if err := c.Bind(&input); err != nil {
		util.RespBadRequest(c, "Invalid format")
		return
	}

	if err := input.Validate(); err != nil {
		util.RespBadRequest(c, err.Error())
		return
	}

	driver := mpg.Driver{
		Name: input.Name,
	}
	if err := h.dbPg.CreateDriver(&driver); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	resp := view.User{
		ID: driver.ID,
	}
	util.RespOK(c, resp)
}

// UpdateDriverLocation Update driver location
func (h *Handler) UpdateDriverLocation(c *gin.Context) {
	var input form.Location
	if err := c.Bind(&input); err != nil {
		util.RespBadRequest(c, "Invalid format")
		return
	}

	// Validate driver is exists in database
	driverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.RespNotFound(c)
		return
	}

	driver, err := h.dbPg.GetDriver(driverID)
	if err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	if driver == nil {
		util.RespNotFound(c)
		return
	}

	// Push to redis
	if err := h.dbRedis.PushDriverLocationGeo(driver.ID, input.Lat, input.Lng); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	// Save to cassandra
	location := mcass.DriverLocation{
		DriverID:  driver.ID,
		Lat:       input.Lat,
		Lng:       input.Lng,
		CreatedAt: time.Now().Unix(),
	}
	if err := h.dbCass.CreateDriverLocation(&location); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	resp := view.DriverLocation{
		ID: driver.ID,
		Location: view.Location{
			Lat: location.Lat,
			Lng: location.Lng,
		},
	}
	fmt.Printf("%+v, %+v\n", input, resp)
	util.RespOK(c, resp)
}

// GetDriverHistory Get driver's history
func (h *Handler) GetDriverHistory(c *gin.Context) {
	// Validate driver is exists in database
	driverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.RespNotFound(c)
		return
	}

	driver, err := h.dbPg.GetDriver(driverID)
	if err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	if driver == nil {
		util.RespNotFound(c)
		return
	}

	// 30 minutes ago
	t := time.Now().Add(-30 * time.Minute)
	history, err := h.dbCass.GetDriverHistory(driver.ID, t)
	if err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	resp := view.PopulateDriverHistory(history)
	util.RespOK(c, resp)
}

// UpdateDriverState Update driver state available/busy
func (h *Handler) UpdateDriverState(c *gin.Context) {
	var input form.DriverState
	if err := c.Bind(&input); err != nil {
		util.RespBadRequest(c, "Invalid format")
		return
	}

	if err := input.Validate(); err != nil {
		util.RespBadRequest(c, err.Error())
		return
	}

	// Get driver
	driverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.RespNotFound(c)
		return
	}

	driver, err := h.dbPg.GetDriver(driverID)
	if err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	if driver == nil {
		util.RespNotFound(c)
		return
	}

	// Update database postgresql
	if err := h.dbPg.UpdateDriverState(driver.ID, input.State); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	// Update driver geo in redis
	switch input.State {
	case mpg.StateAvailable:
		// Get driver latest location
		latestLocation, err := h.dbCass.GetDriverLatestLocation(driver.ID)
		if err != nil {
			util.RespInternalServerError(c, err)
			return
		}

		if err := h.dbRedis.PushDriverLocationGeo(driver.ID, latestLocation.Lat, latestLocation.Lng); err != nil {
			util.RespInternalServerError(c, err)
			return
		}
	case mpg.StateBusy:
		if err := h.dbRedis.RemoveDriverLocationGeo(driver.ID); err != nil {
			util.RespInternalServerError(c, err)
			return
		}

	}

	resp := struct{}{}
	util.RespOK(c, resp)
}
