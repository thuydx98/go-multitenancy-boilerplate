package v1services

import (
	"errors"

	database "go-multitenancy-boilerplate/database"
	helpers "go-multitenancy-boilerplate/helpers"
	models "go-multitenancy-boilerplate/models"
)

type MasterUser models.MasterUser

// Creates a standard user in the database.
// Returns the inserted user id
func CreateMasterUser(email string, password string, accountType int) (uint, error) {

	// Slice for found users.
	var foundUsers []MasterUser

	if err := database.Connection.Select("email").Where("email = ?", email).Find(&foundUsers).Error; err != nil {
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

	var user = MasterUser{Email: email, Password: hash, AccountType: accountType}

	// Run create
	if err := database.Connection.Create(&user).Error; err != nil {
		// Error Handler
		return 0, err
	}

	// Return newly created user ID
	return user.ID, nil
}

// Logs a user in.
func LoginMasterUser(email string, password string) (uint, bool, error) {

	// Create local state user
	var user MasterUser

	// Find the user by email, return error if input is malformed.
	if err := database.Connection.First(&user, "email = ?", email).Error; err != nil {
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
func UpdateMasterUser(id uint, email string, accountType int, firstName string, lastName string, phoneNumber string, recoveryEmail string) (string, error) {

	var user MasterUser

	// Update the basic user information, anything that was set as nil will not be changed.
	if err := database.Connection.Model(&user).Where("id = ?", id).Updates(MasterUser{
		Email:         email,
		AccountType:   accountType,
		FirstName:     firstName,
		LastName:      lastName,
		PhoneNumber:   phoneNumber,
		RecoveryEmail: recoveryEmail,
	}).Error; err != nil {
		return "", err
	}

	return "User Information Successfully Updated.", nil
}

// Deletes a user in the database.
func DeleteMasterUser(id uint) (string, error) {
	var user MasterUser

	if err := database.Connection.Where("id = ?", id).Delete(&user).Error; err != nil {
		return "An error occurred when trying to delete the user", err
	}

	return "The user has been successfully deleted", nil
}

// Get a specific user from the database.
func GetMasterUser(id uint) (*MasterUser, error) {

	var user MasterUser

	if err := database.Connection.Select("id, created_at, updated_at, email, account_type, first_name, last_name, phone_number").Where("id = ? ", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
