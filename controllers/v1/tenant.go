package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	database "go-multitenancy-boilerplate/database"
	middlewares "go-multitenancy-boilerplate/middlewares"
	resources "go-multitenancy-boilerplate/resources/api/v1"
	services "go-multitenancy-boilerplate/services/v1"
)

// Init
func SetupTenantRoutes(router *gin.Engine) {

	users := router.Group("/api/v1/tenants")

	users.Use(middlewares.IfMasterAuthorized(database.Store))
	{
		users.POST("", HandleCreateTenant)
	}
}

// @Summary Attempts to create a new tenant as a privileged user.
// @tags tetants
// @Router api/v1/tenants [Post]
func HandleCreateTenant(c *gin.Context) {

	var json resources.CreateNewTenantRequest

	if err := c.Bind(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "No subdomain identifier was found.")
		return
	}

	outcome, err := services.CreateTenant(json.SubDomainIdentifier)

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, outcome, err.Error())
		return
	}

	resources.Succeeded(c, outcome)
}
