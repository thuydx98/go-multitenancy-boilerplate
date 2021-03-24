package resources

import (
	"fmt"
	helpers "go-multitenancy-boilerplate/helpers"
	res "go-multitenancy-boilerplate/resources/api/v1"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wader/gormstore"
)

type LoginAttempt struct {
	LastLoginAttemptTime time.Time
	LoginAttempts        uint
}

// Checks if a user is logged in with a session to the master dashboard
func HandleMasterLoginAttempt(Store *gormstore.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Try and get a session.
		sessionValues, err := Store.Get(c.Request, "connect.s.id")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		// Check to see if the user is already authorized..
		if sessionValues.ID != "" {

			p := sessionValues.Values["profile"].(HostProfile)

			if p.Authorized == 1 {
				c.JSON(http.StatusOK, gin.H{
					"outcome": "Already Authorized",
					"message": "user already authorized with application.",
				})
				c.Abort()
				return
			}
		}

		// Check our parameters out.
		var json res.LoginRequest

		// Abort if we don't have the correct variables to begin with.
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			fmt.Println("Can't bind request variables for login")
			c.Abort()
			return
		}

		if !helpers.ValidateEmail(json.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			fmt.Println("Email is not in a valid format.")
			c.Abort()
			return
		}

		// Abort if the passed password is not correct.

		// Validate the password being sent.
		if len(json.Password) <= 7 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password was to short, must be longer than 8 characters."})
			c.Abort()
			return
		}

		// Validate the password contains at least one letter and capital
		if !helpers.ContainsCapitalLetter(json.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password does not contain a capital letter."})
			c.Abort()
			return
		}

		// Make sure the password contains at least one special character.
		if !helpers.ContainsSpecialCharacter(json.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The password must contain at least one special character."})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		c.Set("bindedJson", json)

		// Check to see if a new session is found.
		if sessionValues.ID == "" {
			// Setup new session with empty profiles.
			session, err := Store.New(c.Request, "connect.s.id")

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
				c.Abort()
				return
			}

			// Host profile requires little setup
			session.Values["profile"] = newHostProfile()

			// Add the entry record to hostProfile
			session.Values["profile"].(HostProfile).LoginAttempts[json.Email] = &LoginAttempt{LoginAttempts: 1, LastLoginAttemptTime: time.Now().UTC()}

			// Client profile requires no setup
			session.Values["client"] = newClientProfile()

			// Set the session back to the handler for use.
			c.Set("session", session)
			return
		} else {
			// Profile was already found
			h := sessionValues.Values["profile"].(HostProfile)

			// Check if the email used is already in our login attempts.
			loginAttemptsFound, found := h.LoginAttempts[json.Email]

			if !found {
				// email has not been used to login add a new entry
				h.LoginAttempts[json.Email] = &LoginAttempt{LoginAttempts: 1, LastLoginAttemptTime: time.Now().UTC()}

				// Set the session back to the handler for use.
				c.Set("session", sessionValues)
				return
			}

			// Check to see if login attempts exceeds 3 attempts
			if found && loginAttemptsFound.LoginAttempts > 2 {
				// Check to see if last login attempt was over half an hour ago
				if time.Now().Sub(loginAttemptsFound.LastLoginAttemptTime).Minutes() > 30 {
					// reset login attempts to have 2 more.
					loginAttemptsFound.LoginAttempts = 1
					loginAttemptsFound.LastLoginAttemptTime = time.Now().UTC()

					// Set the session back to the handler for use.
					c.Set("session", sessionValues)
					return
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "You have been locked out for too many attempts to login..", "status": "locked out", "timeLeft": 30 - time.Now().Sub(loginAttemptsFound.LastLoginAttemptTime).Minutes()})
					c.Abort()
					return
				}
			}

			if found && loginAttemptsFound.LoginAttempts <= 2 {
				// increase login attempt count
				loginAttemptsFound.LoginAttempts++
				// replace last attempt date
				loginAttemptsFound.LastLoginAttemptTime = time.Now().UTC()

				// Set the session back to the handler for use.
				c.Set("session", sessionValues)
				return
			}
		}
	}
}

// Checks if a user is logged in with a session to the client dashboard
func HandleLoginAttempt(Store *gormstore.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Try and get tenancy identifier
		tenantIdentifier, found := c.Get("tenantIdentifier")

		if !found {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		// Try and get a session.
		sessionValues, err := Store.Get(c.Request, "connect.s.id")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		// Check to see if the user is already authorized..
		if sessionValues.ID != "" {

			p := sessionValues.Values["client"].(ClientProfile)

			authorizationEntry := p.AuthorizationMap[tenantIdentifier.(string)]

			if authorizationEntry == 1 {
				c.JSON(http.StatusOK, gin.H{
					"outcome": "Already Authorized",
					"message": "user already authorized with application.",
				})
				c.Abort()
				return
			}
		}

		// Check our parameters out.
		var json res.LoginRequest

		// Abort if we don't have the correct variables to begin with.
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			fmt.Println("Can't bind request variables for login")
			c.Abort()
			return
		}

		if !helpers.ValidateEmail(json.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			fmt.Println("Email is not in a valid format.")
			c.Abort()
			return
		}

		// Abort if the passed password is not correct.

		// Validate the password being sent.
		if len(json.Password) <= 7 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password was to short, must be longer than 8 characters."})
			c.Abort()
			return
		}

		// Validate the password contains at least one letter and capital
		if !helpers.ContainsCapitalLetter(json.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password does not contain a capital letter."})
			c.Abort()
			return
		}

		// Make sure the password contains at least one special character.
		if !helpers.ContainsSpecialCharacter(json.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The password must contain at least one special character."})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		c.Set("bindedJson", json)

		// Check to see if a new session is found.
		if sessionValues.ID == "" {
			// Setup new session with empty profiles.
			session, err := Store.New(c.Request, "connect.s.id")

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
				c.Abort()
				return
			}

			// Host profile requires little setup
			session.Values["profile"] = newHostProfile()

			// Client profile requires no setup
			session.Values["client"] = newClientProfile()

			session.Values["client"].(ClientProfile).LoginAttempts[tenantIdentifier.(string)] = make(map[string]*LoginAttempt)
			session.Values["client"].(ClientProfile).LoginAttempts[tenantIdentifier.(string)][json.Email] = &LoginAttempt{LoginAttempts: 1, LastLoginAttemptTime: time.Now().UTC()}

			// Set the session back to the handler for use.
			c.Set("session", session)
			return
		} else {
			// Profile was already found
			h := sessionValues.Values["client"].(ClientProfile)

			// Attempt to find tenant entry in login attempts.
			tenantMap, found := h.LoginAttempts[tenantIdentifier.(string)]

			if !found {
				// Create a new entry for the tenant entry in map, also create login attempt
				tenantMap = make(map[string]*LoginAttempt)
				tenantMap[json.Email] = &LoginAttempt{LoginAttempts: 1, LastLoginAttemptTime: time.Now().UTC()}

				// Set the session back to the handler for use.
				c.Set("session", sessionValues)
				return
			}

			// Check if the email used is already in our login attempts.
			loginAttemptsFound, found := tenantMap[json.Email]

			if !found {
				// email has not been used to login add a new entry
				loginAttemptsFound = &LoginAttempt{LoginAttempts: 1, LastLoginAttemptTime: time.Now().UTC()}

				// Set the session back to the handler for use.
				c.Set("session", sessionValues)
				return
			}

			// Check to see if login attempts exceeds 3 attempts
			if found && loginAttemptsFound.LoginAttempts > 2 {
				// Check to see if last login attempt was over half an hour ago
				if time.Now().Sub(loginAttemptsFound.LastLoginAttemptTime).Minutes() > 30 {
					// reset login attempts to have 2 more.
					loginAttemptsFound.LoginAttempts = 1
					loginAttemptsFound.LastLoginAttemptTime = time.Now().UTC()

					// Set the session back to the handler for use.
					c.Set("session", sessionValues)
					return
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "You have been locked out for too many attempts to login..", "status": "locked out", "timeLeft": 30 - time.Now().Sub(loginAttemptsFound.LastLoginAttemptTime).Minutes()})
					c.Abort()
					return
				}
			}

			if found && loginAttemptsFound.LoginAttempts <= 2 {
				// increase login attempt count
				loginAttemptsFound.LoginAttempts++
				// replace last attempt date
				loginAttemptsFound.LastLoginAttemptTime = time.Now().UTC()

				// Set the session back to the handler for use.
				c.Set("session", sessionValues)
				return
			}

		}

	}
}
