package database

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	tenants "go-multitenancy-boilerplate/models/tenants"
	sessions "go-multitenancy-boilerplate/resources/sessions"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/wader/gormstore"
)

var Connection *gorm.DB
var Store *gormstore.Store

func StartDatabaseServices() {

	// Database Connection string
	connectionString := fmt.Sprintf(os.Getenv("CONNECTION_STRING"), os.Getenv("DATABASE_NAME"))
	db, err := gorm.Open(os.Getenv("DIALECT"), connectionString)

	if err != nil {
		fmt.Println(err)
		log.Panic("failed to connect database")
	}

	// Turn logging for the database on.
	db.LogMode(true)

	// Make Master connection available globally.
	Connection = db

	// Now Setup store - Tenant Store
	// Password is passed as byte key method
	Store = gormstore.NewOptions(db, gormstore.Options{
		TableName:       "sessions",
		SkipCreateTable: false,
	}, []byte(os.Getenv("sessionsPassword")))

	// Register session types for consuming in sessions
	gob.Register(sessions.HostProfile{})
	gob.Register(sessions.ClientProfile{})

	// Always attempt to migrate changes to the master tenant schema
	if err := MigrateMasterTenantDatabase(); err != nil {
		fmt.Print("There was an error while trying to migrate the tenant tables..")
		os.Exit(1)
	}

	// attempt to migrate any tenant table changes to all clients.
	AutoMigrateTenantTableChanges()

	// Makes quit Available
	quit := make(chan struct{})

	// Every hour remove dead sessions.
	go Store.PeriodicCleanup(1*time.Hour, quit)
}

// Simply migrates all of the tenant tables
func AutoMigrateTenantTableChanges() {

	var TenantInformation []tenants.TenantConnectionInformation

	Connection.Find(&TenantInformation)

	for _, element := range TenantInformation {

		conn, _ := element.GetConnection()

		if err := MigrateTenantTables(conn); err != nil {
			fmt.Print("An error occurred while attempting to migrate tenant tables")
			os.Exit(1)
		}
	}
}
