package app

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Print(err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")

	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s",
		dbHost, username, dbName, password)

	fmt.Printf("[+] dbUri: %s\n", dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		panic(err)
	}

	db = conn
	db.Debug().AutoMigrate()
}

func GetDB() *gorm.DB {
	return db
}
