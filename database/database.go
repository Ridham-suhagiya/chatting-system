package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func ConnectDatabase(databaseUrl string) (*DB, error) {
	db, error := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if error != nil {
		return nil, error
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func GetDataBaseConnection() (*DB, error) {
	var databaseUrl string = os.Getenv("DATABASE_URL")
	return ConnectDatabase(databaseUrl)
}
