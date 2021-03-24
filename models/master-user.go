package models

// Master User structure
type MasterUser struct {
	Model
	Email         string
	Password      string `json:",omitempty"`
	AccountType   int
	FirstName     string
	LastName      string
	PhoneNumber   string
	RecoveryEmail string
}
