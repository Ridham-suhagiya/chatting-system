package databaseServiceMapper

import (
	"chatting-system-backend/database"
	"chatting-system-backend/service"
	"fmt"
	"os"
)

type DatabaseServiceRepo interface {
	GetService(serviceName string) (interface{}, error)
}

type DatabaseServiceMapper struct {
	services map[string]interface{}
}

func NewServiceMapper() (DatabaseServiceRepo, error) {
	var databaseUrl string = os.Getenv("DATABASE_URL")
	database, err := database.ConnectDatabase(databaseUrl)
	if err != nil {
		fmt.Println("something went wrong in connecting to the database")
		return nil, err
	}
	return &DatabaseServiceMapper{
		services: map[string]interface{}{
			"user": service.NewUserService(database),
		},
	}, nil
}

func (sm *DatabaseServiceMapper) GetService(serviceName string) (interface{}, error) {
	service, exists := sm.services[serviceName]
	if exists {
		return service, nil
	}
	return nil, fmt.Errorf("service not found: %s", serviceName)
}
