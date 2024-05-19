package handler

import (
	"chatting-system-backend/databaseServiceMapper"
	"chatting-system-backend/model"
	"chatting-system-backend/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CreateUserHandler() http.HandlerFunc {
	serviceMapper, err := databaseServiceMapper.NewServiceMapper()
	userSeriviceInterface, err := serviceMapper.GetService("user")
	if err != nil {
		fmt.Println("Something went wrong in fetching sercvice from service interface.")
	}
	userService, ok := userSeriviceInterface.(service.UserService)
	if !ok {
		fmt.Println("Something went wrong in fetching sercvice from service mapper.")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Unable to read request body", http.StatusInternalServerError)
				return
			}
			var user model.User
			fmt.Println(string(body))
			err = json.Unmarshal(body, &user)
			if err != nil {
				response := map[string]interface{}{
					"message": err.Error(),
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
			}
			err = userService.CreateUser(user)
			if err != nil {
				response := map[string]interface{}{
					"message": err.Error(),
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
			}
			w.WriteHeader(http.StatusCreated)
			response := map[string]interface{}{
				"message": "User created successfully",
				"user":    user,
			}
			json.NewEncoder(w).Encode(response)
		}
	}
}

func GetUserHandler() http.HandlerFunc {
	serviceMapper, err := databaseServiceMapper.NewServiceMapper()
	userSeriviceInterface, err := serviceMapper.GetService("user")
	if err != nil {
		fmt.Println("Something went wrong in fetching sercvice from service interface.")
	}
	userService, ok := userSeriviceInterface.(service.UserService)
	if !ok {
		fmt.Println("Something went wrong in fetching sercvice from service mapper.")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		userId := strings.TrimPrefix(r.URL.Path, "/user/")

		fmt.Println(userId, "this is user id")
		user, err := userService.GetUserByID(userId)
		if err != nil {
			response := map[string]interface{}{
				"message": "Record not found",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		response := map[string]interface{}{
			"message": "User fetched successfully",
			"user":    user,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func GetAllUsersHandler() http.HandlerFunc {
	serviceMapper, err := databaseServiceMapper.NewServiceMapper()
	if err != nil {
		fmt.Println("Somethign went wrong in connecting database")
	}
	userServiceInterface, err := serviceMapper.GetService("user")
	if err != nil {
		fmt.Println("Something went wrong in fetching service interface")
	}
	userService, ok := userServiceInterface.(service.UserService)
	if !ok {
		fmt.Println("Something went wrong in asserting to user-service")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}
		w.Header().Set("Content-type", "application/json")
		var response map[string]interface{}
		users, err := userService.GetAllUsers()
		if err != nil {
			response = map[string]interface{}{
				"message": "Something went wrong in fetching the users",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		response = map[string]interface{}{
			"message": "All users fetched Successfully",
			"users":   users,
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
		return
	}

}
