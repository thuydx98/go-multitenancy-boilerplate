package models

// Master User structure
type MasterUser struct {
	Model
	Email         string `json:"email"`
	Password      string `json:",omitempty"`
	AccountType   int    `json:"account_type"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	PhoneNumber   string `json:"phone_number"`
	RecoveryEmail string `json:"recovery_email"`
}
