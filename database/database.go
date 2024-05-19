package database

import (
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
