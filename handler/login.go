package handler

import (
	"chatting-system-backend/databaseServiceMapper"
	"chatting-system-backend/objectTypes"
	"chatting-system-backend/service"
	"chatting-system-backend/utils"
	"fmt"
	"net/http"
	"time"
)

func LoginHandler() http.HandlerFunc {
	serviceMapper, err := databaseServiceMapper.NewServiceMapper()
	if err != nil {
		fmt.Println("Something went wrong in fetching service")
	}

	serviceInterface, err := serviceMapper.GetService("user")
	if err != nil {
		fmt.Println("Something went wrong in fetching service interface")
	}
	userService := serviceInterface.(service.UserService)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {

			headers := map[string]interface{}{
				"statusCode":  http.StatusNotFound,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "Method not allowed",
			}

			utils.WriteIntoTheResponse(w, params)
			return
		}
		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			fmt.Println(err.Error(), "this is the error")
			headers := map[string]interface{}{
				"statusCode":  http.StatusInternalServerError,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "Something went wrong while handling unmarshal",
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := userService.GetUserByEmailId(email)
		fmt.Println("this is not an erorr", user, err)
		if err != nil {
			headers := map[string]interface{}{
				"statusCode":  http.StatusNotFound,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "User Credentials invalid.",
			}
			utils.WriteIntoTheResponse(w, params)
			fmt.Println(w)
			return
		}

		if !userService.ValidatePassword(user, password) {
			fmt.Println("this is not an erorr", user, err, password)
			headers := map[string]interface{}{
				"statusCode":  http.StatusNotFound,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "Invalid user credentials",
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}
		userDetails := &objectTypes.LoginCredentials{
			User: user,
		}

		// Set session cookie
		token, err := utils.GenerateJWT(userDetails)

		if err != nil {
			headers := map[string]interface{}{
				"statusCode":  http.StatusInternalServerError,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "Something went wrong in generating token",
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}

		headers := map[string]interface{}{
			"statusCode":  http.StatusOK,
			"auth_token":  token,
			"contentType": "application/json",
		}
		params := utils.ResponseParams{
			Header:  headers,
			Message: "User Credentials valid.",
			Details: userDetails,
		}
		utils.WriteIntoTheResponse(w, params)
	}
}

func LogOutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			headers := map[string]interface{}{
				"statusCode":  http.StatusNotFound,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "Method not allowed",
			}

			utils.WriteIntoTheResponse(w, params)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "usert",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Use true if serving over HTTPS
			SameSite: http.SameSiteLaxMode,
		})
		headers := map[string]interface{}{
			"statusCode":  http.StatusOK,
			"contentType": "application/json",
		}
		params := utils.ResponseParams{
			Header:  headers,
			Message: "User Logged out successfully",
		}

		utils.WriteIntoTheResponse(w, params)
	}
}

func CheckUserLoggedIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {

			headers := map[string]interface{}{
				"statusCode":  http.StatusNotFound,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "Method not allowed",
			}

			utils.WriteIntoTheResponse(w, params)
			return
		}
		authToken, err := utils.ValidateJWT(r.Header.Get("auth_token"))
		if err != nil {
			headers := map[string]interface{}{
				"statusCode":  http.StatusUnauthorized,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "Invalid or expired token",
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}
		headers := map[string]interface{}{
			"statusCode":  http.StatusOK,
			"contentType": "application/json",
			"auth_token":  authToken,
		}
		params := utils.ResponseParams{
			Header:  headers,
			Message: "User is logged in",
		}
		utils.WriteIntoTheResponse(w, params)

	}
}
