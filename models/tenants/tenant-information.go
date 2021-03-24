package models

import (
	"strings"

	"go-multitenancy-boilerplate/models"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type TenantConnectionInformation struct {
	models.Model
	TenantId                  uint `gorm:"AUTO_INCREMENT"`
	TenantSubDomainIdentifier string
	ConnectionString          string
}

// Helper method that create and returns the database connection.
func (t TenantConnectionInformation) GetConnection() (*gorm.DB, error) {

	if len(strings.TrimSpace(t.ConnectionString)) == 0 {
		return nil, errors.New("Connection string was not found or was empty..")
	}

	db, err := gorm.Open("postgres", t.ConnectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
