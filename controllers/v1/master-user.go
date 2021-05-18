package v1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	database "go-multitenancy-boilerplate/database"
	helpers "go-multitenancy-boilerplate/helpers"
	middlewares "go-multitenancy-boilerplate/middlewares"
	resources "go-multitenancy-boilerplate/resources/api/v1"
	ss "go-multitenancy-boilerplate/resources/sessions"
	services "go-multitenancy-boilerplate/services/v1"
)

// Init
func SetupMasterUserRoutes(router *gin.Engine) {

	users := router.Group("/api/v1/master/users")

	users.POST("login", ss.HandleMasterLoginAttempt(database.Store), HandleMasterLogin)

	users.Use(middlewares.IfMasterAuthorized(database.Store))
	{
		// POST
		users.POST("", HandleMasterCreateUser)
		users.POST("logout", HandleMasterLogout)

		// PUT
		users.PUT("", HandleMasterUpdateUserDetails)

		// GET
		users.GET("{id}", HandleMasterGetUserById)
		users.GET("me", HandleMasterGetCurrentUser)

		// DELETE
		users.DELETE("", HandleMasterDeleteUser)
	}
}

// @Summary Create a new user
// @tags master/users
// @Router /master/api/users/create [post]
func HandleMasterCreateUser(c *gin.Context) {

	// Binds Model and handles validation.
	var json resources.CreateUserRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "Incorrect details supplied, please try again.")
		return
	}

	if !helpers.ValidateEmail(json.Email) {
		resources.Failed(c, http.StatusBadRequest, "Email is incorrect, please try again.")
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

	// Attempt to create a user.
	insertedId, err := services.CreateMasterUser(json.Email, json.Password, json.Type)

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	resources.Succeeded(c, insertedId)
}

// @Summary Attempt to login using user details
// @tags master/users
// @Router /master/api/users/login [post]
func HandleMasterLogin(c *gin.Context) {

	bindJson, _ := c.Get("bindedJson")

	json := bindJson.(resources.LoginRequest)

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

	// Get our session from database.
	session, exists := c.Get("session")

	if !exists {
		resources.Failed(c, http.StatusUnprocessableEntity, "Something went wrong while trying to process that, please try again.")
		return
	}

	userId, outcome, err := services.LoginMasterUser(json.Email, json.Password)

	if err != nil {

		// Save changes to our session if an error occurred and we need to abort early..
		if err := database.Store.Save(c.Request, c.Writer, session.(*sessions.Session)); err != nil {
			fmt.Print(err)
		}

		// Were sending 422 as there is a validation concern.
		resources.Failed(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Create a copy of the host profile
	hostProfile := session.(*sessions.Session).Values["profile"].(ss.HostProfile)

	// Set session values to authorized
	hostProfile.Authorized = 1
	hostProfile.AuthorizedTime = time.Now().UTC()
	hostProfile.UserId = userId

	// Reset login attempts once successfully logged in.
	hostProfile.LoginAttempts[json.Email].LoginAttempts = 0

	// Set host profile back to values.
	session.(*sessions.Session).Values["profile"] = hostProfile

	// Save changes to our session.
	if err := database.Store.Save(c.Request, c.Writer, session.(*sessions.Session)); err != nil {
		fmt.Print(err)
	}

	resources.Succeeded(c, outcome)
}

// @Summary Logs a user out of the system
// @tags master/users
// @Router /master/api/users/logout [post]
func HandleMasterLogout(c *gin.Context) {

	// Binds Model and handles validation.
	var json resources.CreateUserRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "Incorrect details supplied, please try again.")
		return
	}

	// Get our session from database.
	session, err := database.Store.Get(c.Request, "connect.s.id")

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Create a copy of the host profile
	hostProfile := session.Values["profile"].(ss.HostProfile)

	// Set session values to unauthorized
	hostProfile.Authorized = 0

	// Set host profile back to values.
	session.Values["profile"] = hostProfile

	// Save changes to our session.
	if err := database.Store.Save(c.Request, c.Writer, session); err != nil {
		fmt.Print(err)
	}

	resources.Succeeded(c, "You have successfully logged out of your account.")
}

// @Summary Updates a users details
// @tags master/users
// @Router /master/api/users/updateUserDetails [post]
func HandleMasterUpdateUserDetails(c *gin.Context) {
	var json resources.UpdateUserRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "Missing required fields, please try again.")
		log.Println(err)
		return
	}

	outcome, err := services.UpdateMasterUser(json.Id, json.Email, json.AccountType, json.FirstName, json.LastName, json.PhoneNumber, json.RecoveryEmail)

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resources.Succeeded(c, outcome)
}

// @Summary Deletes a user using a user id
// @tags master/users
// @Router /master/api/users/deleteUser [delete]
func HandleMasterDeleteUser(c *gin.Context) {
	var json resources.DeleteUserRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "Missing required fields, please try again.")
		return
	}

	outcome, err := services.DeleteMasterUser(json.Id)

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resources.Succeeded(c, outcome)
}

// @Summary Attempts to get a existing user by id
// @tags master/users
// @Router /master/api/users/getUserById [get]
func HandleMasterGetUserById(c *gin.Context) {
	// Were using delete params as it shares the same interface.
	var json resources.DeleteUserRequest

	if err := c.Bind(&json); err != nil {
		resources.Failed(c, http.StatusBadRequest, "No user ID found, please try again.")
		return
	}

	outcome, err := services.GetMasterUser(json.Id)

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resources.Succeeded(c, outcome)
}

// @Summary Attempts to get the currently logged in user using there session id.
// @tags master/users
// @Router /master/api/users/getCurrentUser [get]
func HandleMasterGetCurrentUser(c *gin.Context) {

	// Get the currently logged int user id.
	userId := c.MustGet("userId")

	outcome, err := services.GetMasterUser(userId.(uint))

	if err != nil {
		resources.Failed(c, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resources.Succeeded(c, outcome)
}
