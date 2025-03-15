package handler

import (
	"chatting-system-backend/databaseServiceMapper"
	"chatting-system-backend/model"
	"chatting-system-backend/service"
	"chatting-system-backend/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateLink() http.HandlerFunc {
	serviceMapper, err := databaseServiceMapper.NewServiceMapper()
	if err != nil {
		fmt.Println("Something went wrong in fetching service")
	}
	linkServiceInterface, err := serviceMapper.GetService("link")
	if err != nil {
		fmt.Println("Something went wrong in fetching service")
	}

	linkService, ok := linkServiceInterface.(service.LinkService)
	if !ok {
		fmt.Println("Something went wrong in fetching service interface")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			err := map[string]interface{}{
				"message": "Method not allowed",
				"name":    "MethodNotAllowed",
			}
			fmt.Println(err)
			headers := map[string]interface{}{

				"statusCode":  http.StatusNotFound,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Details: err,
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			response := map[string]interface{}{
				"message": err.Error(),
			}
			headers := map[string]interface{}{
				"statusCode":  http.StatusUnprocessableEntity,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Details: response,
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}
		userIdStr := r.FormValue("user_id")
		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			response := map[string]interface{}{
				"message": "invalid user",
			}
			headers := map[string]interface{}{
				"statusCode":  http.StatusUnauthorized,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Details: response,
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}

		expiryDateStr := r.FormValue("expiry_date")
		if expiryDateStr == "" {
			expiryDateStr = "2026-12-31"
		}
		expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
		if err != nil {
			response := map[string]interface{}{
				"message": "invalid expiry date",
			}
			headers := map[string]interface{}{
				"statusCode":  http.StatusBadRequest,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Details: response,
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}
		link := &model.ChatLinks{
			UserId:   userId,
			ExpiryAt: expiryDate,
		}

		err = linkService.CreateLink(link)
		if err != nil {
			fmt.Print(err)
			headers := map[string]interface{}{
				"statusCode":  http.StatusUnprocessableEntity,
				"contentType": "application/json",
			}
			params := utils.ResponseParams{
				Header:  headers,
				Message: "Link not created.",
			}
			utils.WriteIntoTheResponse(w, params)
			return
		}
		headers := map[string]interface{}{
			"statusCode":  http.StatusCreated,
			"contentType": "application/json",
		}
		params := utils.ResponseParams{
			Header:  headers,
			Message: "Link created successfully",
			Details: link,
		}
		utils.WriteIntoTheResponse(w, params)
	}
}
