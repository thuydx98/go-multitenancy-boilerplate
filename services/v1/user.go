package v1services

import (
	"errors"
	helpers "go-multitenancy-boilerplate/helpers"
	"go-multitenancy-boilerplate/models"

	"github.com/jinzhu/gorm"
)

// Creates a standard user in the database.
// Returns the inserted user id
func CreateUser(email string, password string, accountType int, connection *gorm.DB) (uint, error) {

	// Slice for found users.
	var foundUsers []models.User

	if err := connection.Select("email").Where("email = ?", email).Find(&foundUsers).Error; err != nil {
		return 0, err
	}

	// If duplicate email address has been found return.
	if len(foundUsers) > 0 {
		return 0, errors.New("A user with that email address already exists")
	}

	// Hash the password so it's not clear text.
	// Run outside of the if statements so we can grab the result outside of local scope.
	hash, hashErr := helpers.HashPassword([]byte(password))

	if hashErr != nil {
		return 0, hashErr
	}

	var user = models.User{Email: email, Password: hash, AccountType: accountType}

	// Run create
	if err := connection.Create(&user).Error; err != nil {
		// Error Handler
		return 0, err
	}

	// Return newly created user ID
	return user.ID, nil
}

// Logs a user in.
func LoginUser(email string, password string, connection *gorm.DB) (uint, bool, error) {

	// Create local state user
	var user models.User

	// Find the user by email, return error if input is malformed.
	if err := connection.First(&user, "email = ?", email).Error; err != nil {
		return 0, false, err
	}

	// Now we've found a user send off the hashed password and sent password for decoding.
	if result := helpers.CheckPasswordHash(password, user.Password); result != true {
		// Passwords do not match
		return 0, false, errors.New("passwords did not match")
	}

	// Checks have bee passed return true
	return user.ID, true, nil
}

// Updates a user in the database.
// A separate method is called when updating a company id
func UpdateUser(id uint, email string, accountType int, firstName string, lastName string, phoneNumber string, recoveryEmail string, connection *gorm.DB) (string, error) {

	var user models.User

	// Update the basic user information, anything that was set as nil will not be changed.
	err := connection.Model(&user).Where("id = ?", id).Updates(models.User{
		Email:         email,
		AccountType:   accountType,
		FirstName:     firstName,
		LastName:      lastName,
		PhoneNumber:   phoneNumber,
		RecoveryEmail: recoveryEmail,
	}).Error

	if err != nil {
		return "", err
	}

	return "User Information Successfully Updated", nil
}

// Deletes a user in the database.
func DeleteUser(id uint, connection *gorm.DB) (string, error) {
	var user models.User

	if err := connection.Where("id = ?", id).Delete(&user).Error; err != nil {
		return "An error occurred when trying to delete the user", err
	}

	return "The user has been successfully deleted", nil
}

// Get a specific user from the database.
func GetUser(id uint, connection *gorm.DB) (*models.User, error) {

	var user models.User

	if err := connection.Select("id, created_at, updated_at, deleted_at, email, account_type, company_id, first_name, last_name").Where("id = ? ", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
