package v1resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response struct for return
type Response struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	Errors  []interface{} `json:"error,omitempty"`
}

//Message returns map data
func Failed(c *gin.Context, status int, message string, errors ...interface{}) {
	c.AbortWithStatusJSON(status, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func Succeeded(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}
