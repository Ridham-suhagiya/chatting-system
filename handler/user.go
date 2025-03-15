package handler

import (
	"chatting-system-backend/databaseServiceMapper"
	"chatting-system-backend/model"
	"chatting-system-backend/service"
	"chatting-system-backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func CreateUserHandler() http.HandlerFunc {
	serviceMapper, _ := databaseServiceMapper.NewServiceMapper()
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
			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				response := map[string]interface{}{
					"message": err.Error(),
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
			}
			email := r.FormValue("email")
			password := r.FormValue("password")
			username := r.FormValue("username")

			user := &model.User{
				Username: username,
				Email:    email,
				Password: password,
			}

			err = userService.CreateUser(user)
			if err != nil {
				fmt.Print(err)
				headers := map[string]interface{}{
					"statusCode":  http.StatusUnprocessableEntity,
					"contentType": "application/json",
				}
				params := utils.ResponseParams{
					Header:  headers,
					Message: "User not created.",
				}
				utils.WriteIntoTheResponse(w, params)
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
	serviceMapper, _ := databaseServiceMapper.NewServiceMapper()
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
	}

}
