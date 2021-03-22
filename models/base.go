package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *gorm.DB

//Model is sample of common table structure
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	// username := os.Getenv("db_user")
	// password := os.Getenv("db_pass")
	// dbName := os.Getenv("db_name")
	// dbHost := os.Getenv("db_host")
	// dbPort := os.Getenv("db_port")

	// sql := mysql.Config{}
	// log.Println(sql)

	// conn, err := gorm.Open("mysql", username+":"+password+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=True&loc=Asia%2FKolkata")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"123456",
		"go-boilerplate",
	)

	fmt.Println("DB Connecting", connectionString)
	dbCon, err := gorm.Open("postgres", connectionString)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("DB Connected")
	db = dbCon
	db.SingularTable(true)
	db.LogMode(true)

	db.Debug().AutoMigrate(
		&User{},
	)
}

// GetDB function return the instance of db
func GetDB() *gorm.DB {
	return db
}
