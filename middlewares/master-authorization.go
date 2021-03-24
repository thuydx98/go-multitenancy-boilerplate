package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wader/gormstore"

	resources "go-multitenancy-boilerplate/resources/api/v1"
	ss "go-multitenancy-boilerplate/resources/sessions"
)

// Checks if a user is logged in with a session to the master dashboard;
func IfMasterAuthorized(Store *gormstore.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionValues, err := Store.Get(c.Request, "connect.s.id")
		if err != nil {
			resources.Failed(c, http.StatusUnauthorized, "You are not authorized to view this.")
			return
		}

		profile := sessionValues.Values["profile"].(ss.HostProfile)
		if profile.Authorized != 1 {
			resources.Failed(c, http.StatusUnauthorized, "You are not authorized to view this.")
			return
		}

		// Pass the user id into the handler.
		c.Set("userId", sessionValues.Values["userId"])
	}
}
