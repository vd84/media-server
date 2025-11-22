package db

import (
	"fmt"
	"log"
	"mediaserver/data"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	port     = 1433
	database = "media_server"
)

func ConnectToLocalDb() (*gorm.DB, error) {
	dsn := "host=localhost user=user1234 password=user1234 dbname=localdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToDb() (*gorm.DB, error) {
	server := os.Getenv("DB_SERVER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	db, err := gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	fmt.Printf("Connected!")
	return db, err
}

func RunMigration(db *gorm.DB) {
	if err := db.AutoMigrate(&data.User{}, &data.Movie{}, &data.OmdbMovie{}, &data.Rating{}, &data.Episode{}, &data.Series{}); err != nil {
		log.Fatal("Migration failed:", err)
	}
}
