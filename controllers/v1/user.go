package v1

import (
	"encoding/json"
	res "go-multitenancy-boilerplate/resources/api/v1"
	services "go-multitenancy-boilerplate/services/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserList function will give you the list of users
func UserList(c *gin.Context) {
	var userService services.UserService

	// decode the request body into struct and failed if any error occur
	err := json.NewDecoder(c.Request.Body).Decode(&userService.User)
	if err != nil {
		res.Respond(c.Writer, res.Message(http.StatusBadRequest, "Invalid request"))
		return
	}

	// call service
	resp := userService.UserList()

	res.Respond(c.Writer, resp)
}
