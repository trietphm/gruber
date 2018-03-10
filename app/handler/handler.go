package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trietphm/gruber/app/form"
	"github.com/trietphm/gruber/app/view"
	"github.com/trietphm/gruber/database"
	"github.com/trietphm/gruber/model"
	"github.com/trietphm/gruber/util"
)

// Handler Manager handler functions
type Handler struct {
	DB database.Database
}

// NewEngine Setup API router
func NewEngine(db database.Database) (*gin.Engine, error) {
	engine := gin.Default()

	handler := Handler{
		DB: db,
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

	passenger := model.Passenger{
		Name: input.Name,
	}
	if err := h.DB.CreatePassenger(&passenger); err != nil {
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

	drivers, err := h.DB.GetNearestDrivers(input.Location.Lat, input.Location.Lng)
	if err != nil {
		util.RespInternalServerError(c, err)
		return
	}

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

	driver := model.Driver{
		Name: input.Name,
	}
	if err := h.DB.CreateDriver(&driver); err != nil {
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

	driver, err := h.DB.GetDriver(driverID)
	if err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	if driver == nil {
		util.RespNotFound(c)
		return
	}

	// Create location
	driverLocation := model.DriverLocation{
		DriverID: driver.ID,
		Lat:      input.Lat,
		Lng:      input.Lng,
	}

	if err := h.DB.CreateDriverLocation(&driverLocation); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	resp := view.DriverLocation{
		ID: driverLocation.ID,
		Location: view.Location{
			Lat: driverLocation.Lat,
			Lng: driverLocation.Lng,
		},
	}
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

	driver, err := h.DB.GetDriver(driverID)
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
	history, err := h.DB.GetDriverHistory(driver.ID, t)
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

	driver, err := h.DB.GetDriver(driverID)
	if err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	if driver == nil {
		util.RespNotFound(c)
		return
	}

	if err := h.DB.UpdateDriverState(driver.ID, input.State); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	resp := struct{}{}
	util.RespOK(c, resp)
}
