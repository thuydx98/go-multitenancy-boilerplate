package routers

import (
	v1 "go-multitenancy-boilerplate/controllers/v1"

	"github.com/gin-gonic/gin"
)

// SetupRouter function will perform all route operations
func SetupRouter() *gin.Engine {

	router := gin.Default()

	// Giving access to storage folder
	router.Static("/storage", "storage")

	// Giving access to template folder
	router.Static("/templates", "templates")
	router.LoadHTMLGlob("templates/*")

	router.Use(CORSMiddleware())

	// API route for version 1
	v1.SetupUserRoutes(router)
	v1.SetupMasterUserRoutes(router)
	v1.SetupTenantRoutes(router)

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
