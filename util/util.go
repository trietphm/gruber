package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespBadRequest Response HTTP status Bad Request error with a Json message `{"message":<message>}`
func RespBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, map[string]string{"message": message})
}

// RespInternalServerError Response HTTP status Internal server error message and log the error
func RespInternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, map[string]string{"message": "INTERNAL SERVER ERROR"})
	//TODO Log error instead of panic
	panic(err)
}

// RespOK Response HTTP status OK with json data
func RespOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// RespNotFound Response HTTP status Not found with empty data
func RespNotFound(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotFound)
}
