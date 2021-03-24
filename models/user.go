package models

//User structure
type User struct {
	Model
	Email         string `gorm:"type:varchar(50)" json:"email" validate:"required,email"`
	Password      string `json:",omitempty"`
	AccountType   int
	FirstName     string `gorm:"type:varchar(50)" json:"first_name"`
	LastName      string `gorm:"type:varchar(50)" json:"last_name"`
	PhoneNumber   string
	RecoveryEmail string
}
