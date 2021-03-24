package v1resources

//UserResponse struct
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty" validate:"required,email"`
}

type CreateUserRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Type     int    `form:"type" json:"type"`
}

type UpdateUserRequest struct {
	Id            uint   `form:"id" json:"id" binding:"required"`
	Email         string `form:"email" json:"email"`
	AccountType   int    `form:"accountType" json:"accountType"`
	FirstName     string `form:"firstName" json:"firstName"`
	LastName      string `form:"lastName" json:"lastName"`
	PhoneNumber   string `form:"phoneNumber" json:"phoneNumber"`
	RecoveryEmail string `form:"recoveryEmail" json:"recoveryEmail"`
}

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type DeleteUserRequest struct {
	Id uint `form:"id" json:"id" binding:"required"`
}
