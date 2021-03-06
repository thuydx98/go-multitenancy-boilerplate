package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"

	database "go-multitenancy-boilerplate/database"
	helpers "go-multitenancy-boilerplate/helpers"
	middlewares "go-multitenancy-boilerplate/middlewares"
	resources "go-multitenancy-boilerplate/resources/api/v1"
	services "go-multitenancy-boilerplate/services/v1"
)

// Init
func SetupUserRoutes(router *gin.Engine) {

	users := router.Group("/api/v1/users")

	// Un-authorize APIs
	users.Use(middlewares.FindTenancy(database.Connection))
	{
		users.POST("login", HandleLogin)

		// Authorized APIs
		users.Use(middlewares.IfAuthorized(database.Store))
		{
			users.GET("{id}", HandleGetUserById)
			users.GET("me", HandleGetCurrentUser)

			users.POST("", HandleCreateUser)

			users.PUT("", HandleUpdateUserDetails)

			users.DELETE("", HandleDeleteUser)
		}
	}
}

// @Summary Create a new user
// @tags users
// @Router /api/users/create [post]
func HandleCreateUser(c *gin.Context) {

	// Binds Model and handles validation.
	var json resources.CreateUserRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "Incorrect details supplied, please try again.")
		return
	}

	if !helpers.ValidateEmail(json.Email) {
		resources.Failed(c, http.StatusBadRequest, "Email or Password provided are incorrect, please try again.")
		return
	}

	// Validate the password being sent.
	if len(json.Password) <= 7 {
		resources.Failed(c, http.StatusBadRequest, "The specified password was to short, must be longer than 8 characters.")
		return
	}

	// Validate the password contains at least one letter and capital
	if !helpers.ContainsCapitalLetter(json.Password) {
		resources.Failed(c, http.StatusBadRequest, "The specified password does not contain a capital letter.")
		return
	}

	// Make sure the password contains at least one special character.
	if !helpers.ContainsSpecialCharacter(json.Password) {
		resources.Failed(c, http.StatusBadRequest, "The password must contain at least one special character.")
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	// Attempt to create a user.
	insertedId, err := services.CreateUser(json.Email, json.Password, json.Type, db.(*gorm.DB))

	if err != nil {
		resources.Failed(c, http.StatusBadRequest, "Something went wrong while trying to process that, please try again.", err.Error())
		return
	}

	resources.Succeeded(c, gin.H{
		"id": insertedId,
	})
}

// @Summary Attempt to login using user details
// @tags users
// @Router /api/users/login [post]
func HandleLogin(c *gin.Context) {

	var json resources.LoginRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "Missing required fields, please try again.")
		return
	}

	if !helpers.ValidateEmail(json.Email) {
		resources.Failed(c, http.StatusBadRequest, "Email or Password provided are incorrect, please try again.")
		return
	}

	// Validate the password being sent.
	if len(json.Password) <= 7 {
		resources.Failed(c, http.StatusBadRequest, "The specified password was to short, must be longer than 8 characters.")
		return
	}

	// Validate the password contains at least one letter and capital
	if !helpers.ContainsCapitalLetter(json.Password) {
		resources.Failed(c, http.StatusBadRequest, "The specified password does not contain a capital letter.")
		return
	}

	// Make sure the password contains at least one special character.
	if !helpers.ContainsSpecialCharacter(json.Password) {
		resources.Failed(c, http.StatusBadRequest, "The password must contain at least one special character.")
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	session, exists := c.Get("session")

	if !exists {
		resources.Failed(c, http.StatusUnprocessableEntity, "Something went wrong while trying to process that, please try again.")
		return
	}

	userId, _, err := services.LoginUser(json.Email, json.Password, db.(*gorm.DB))

	if err != nil {

		// Save changes to our session if an error occurred and we need to abort early..
		if err := database.Store.Save(c.Request, c.Writer, session.(*sessions.Session)); err != nil {
			fmt.Print(err)
		}

		// Were sending 422 as there is a validation concern.
		resources.Failed(c, http.StatusUnprocessableEntity, "Something went wrong while trying to process that, please try again.", err.Error())
		return
	}

	// @todo make this into a map so userid's can be multiple.
	session.(*sessions.Session).Values["userId"] = userId

	if err := database.Store.Save(c.Request, c.Writer, session.(*sessions.Session)); err != nil {
		fmt.Print(err)
	}

	resources.Succeeded(c, "You have successfully logged into your account.")
}

// @Summary Updates a users details
// @tags users
// @Router /api/users/updateUserDetails [post]
func HandleUpdateUserDetails(c *gin.Context) {
	var json resources.UpdateUserRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "Missing required fields, please try again.")
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	outcome, err := services.UpdateUser(json.Id, json.Email, json.AccountType, json.FirstName, json.LastName, json.PhoneNumber, json.RecoveryEmail, db.(*gorm.DB))

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, "Something went wrong while trying to process that, please try again.")
		return
	}

	resources.Succeeded(c, outcome)
}

// @Summary Deletes a user using a user id
// @tags users
// @Router /api/users/deleteUser [delete]
func HandleDeleteUser(c *gin.Context) {

	var json resources.DeleteUserRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "Missing required fields, please try again.")
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	outcome, err := services.DeleteUser(json.Id, db.(*gorm.DB))

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, "Something went wrong while trying to process that, please try again.")
		log.Println(err)
		return
	}

	resources.Succeeded(c, outcome)
}

// @Summary Attempts to get a existing user by id
// @tags users
// @Router /api/users/getUserById [get]
func HandleGetUserById(c *gin.Context) {
	// Were using delete params as it shares the same interface.
	var json resources.DeleteUserRequest

	if err := c.Bind(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "No user ID found, please try again.")
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	outcome, err := services.GetUser(json.Id, db.(*gorm.DB))

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, "Something went wrong while trying to process that, please try again.", err.Error())
		log.Println(err)
		return
	}

	resources.Succeeded(c, outcome)
}

// @Summary Attempts to get the currently logged in user using there session id.
// @tags users
// @Router /api/users/getCurrentUser [get]
func HandleGetCurrentUser(c *gin.Context) {

	// Get the currently logged int user id.
	userId := c.MustGet("userId")

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	outcome, err := services.GetUser(userId.(uint), db.(*gorm.DB))

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, "Something went wrong while trying to process that, please try again.", err.Error())
		log.Println(err)
		return
	}

	resources.Succeeded(c, outcome)
}
