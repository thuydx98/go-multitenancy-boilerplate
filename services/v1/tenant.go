package v1services

import (
	"fmt"
	"os"
	"strings"

	database "go-multitenancy-boilerplate/database"
	tenants "go-multitenancy-boilerplate/models/tenants"
)

// Create a tenant using a domain identifier
func CreateTenant(subDomainIdentifier string) (msg string, err error) {

	// Create new database to hold client.
	databaseName := strings.ToLower(subDomainIdentifier) + os.Getenv("SUFFIX_TENANT_DATABASE_NAME")
	connectionString := fmt.Sprintf(os.Getenv("CONNECTION_STRING"), databaseName)

	if err := database.Connection.Exec("CREATE DATABASE \"" + databaseName + "\" OWNER postgres").Error; err != nil {
		return "error making the database", err
	}

	tenant := tenants.TenantConnectionInformation{
		TenantSubDomainIdentifier: subDomainIdentifier,
		ConnectionString:          connectionString,
	}

	if err := database.Connection.Create(&tenant).Error; err != nil {
		return "error inserting the new database record", err
	}

	tenConn, tenConErr := tenant.GetConnection()

	if tenConErr != nil {
		return "error creating the connection using connection method", err
	}

	if migrateErr := database.MigrateTenantTables(tenConn); migrateErr != nil {
		return "error attempting to migrate the existing tables to new database", migrateErr
	}

	return "New Tenant has been successfully made", nil
}
