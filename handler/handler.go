package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trietphm/gruber/database"
	"github.com/trietphm/gruber/form"
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
	router.POST("/passenger", handler.CreatePassenger)
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

	if input.Name == "" {
		util.RespBadRequest(c, "Name can not be empty")
		return
	}

	passenger := model.Passenger{
		Name: input.Name,
	}
	if err := h.DB.CreatePassenger(&passenger); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	resp := struct {
		ID int `json:"id"`
	}{
		ID: passenger.ID,
	}
	util.RespOK(c, resp)
}

// RequestDrivers Get list nearest drivers
func (h *Handler) RequestDrivers(c *gin.Context) {

}

// CreateDriver Sign up driver
func (h *Handler) CreateDriver(c *gin.Context) {
	var input form.Driver
	if err := c.Bind(&input); err != nil {
		util.RespBadRequest(c, "Invalid format")
		return
	}

	if input.Name == "" {
		util.RespBadRequest(c, "Name can not be empty")
		return
	}

	driver := model.Driver{
		Name: input.Name,
	}
	if err := h.DB.CreateDriver(&driver); err != nil {
		util.RespInternalServerError(c, err)
		return
	}

	resp := struct {
		ID int `json:"id"`
	}{
		ID: driver.ID,
	}
	util.RespOK(c, resp)
}

// UpdateDriverLocation Update driver location
func (h *Handler) UpdateDriverLocation(c *gin.Context) {

}

// GetDriverHistory Get driver's history
func (h *Handler) GetDriverHistory(c *gin.Context) {

}

// UpdateDriverState Update driver state available/busy
func (h *Handler) UpdateDriverState(c *gin.Context) {
	var input form.DriverState
	if err := c.Bind(&input); err != nil {
		util.RespBadRequest(c, "Invalid format")
		return
	}

	// Validate
	if input.State == "" || (input.State != model.StateAvailable && input.State != model.StateBusy) {
		util.RespBadRequest(c, "Invalid state")
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
