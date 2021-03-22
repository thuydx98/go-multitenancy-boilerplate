package v1service

import (
	"go-multitenancy-boilerplate/models"
	res "go-multitenancy-boilerplate/resources/api/v1"
)

type UserService struct {
	User models.User
}

// UserList function returns the list of users
func (us *UserService) UserList() map[string]interface{} {
	user := us.User
	userData := res.UserResponse{
		ID:    user.ID,
		Name:  "test",
		Email: "test@gmail.com",
	}

	response := res.Message(0, "This is from version 1 api")
	response["data"] = userData
	return response
}
