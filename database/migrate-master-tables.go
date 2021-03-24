package database

import (
	models "go-multitenancy-boilerplate/models"
	tenants "go-multitenancy-boilerplate/models/tenants"
)

/**
This method uses the base tenant connection set out within init.
*/
func MigrateMasterTenantDatabase() error {

	if err := Connection.AutoMigrate(&tenants.TenantConnectionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&tenants.TenantSubscriptionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&tenants.TenantSubscriptionType{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&models.MasterUser{}).Error; err != nil {
		return err
	}

	return nil
}
