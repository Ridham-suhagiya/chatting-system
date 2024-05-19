package main

import (
	"chatting-system-backend/handler"
	"fmt"
	"net/http"
)

func main() {
	// Define routes and their corresponding handler functions
	http.HandleFunc("/user", handler.CreateUserHandler())
	http.HandleFunc("/user/", handler.GetUserHandler())
	http.HandleFunc("/users", handler.GetAllUsersHandler())
	// Start the web server
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
