package main

import (
	"chatting-system-backend/handler"
	"chatting-system-backend/middleware"
	"fmt"
	"net/http"
)

func ApplyMiddleWare(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h.ServeHTTP
}

func main() {

	mux := http.NewServeMux()

	// Web Socket
	go mux.Handle("/ws", handler.WebSocketHandler())

	// User routes
	mux.Handle("/user", ApplyMiddleWare(handler.CreateUserHandler(), middleware.CORSMiddleware))
	mux.Handle("/user/", ApplyMiddleWare(handler.GetUserHandler(), middleware.CheckSession))
	mux.Handle("/users", ApplyMiddleWare(handler.GetAllUsersHandler(), middleware.CheckSession))

	// Link routes
	mux.Handle("/link", ApplyMiddleWare(handler.CreateLink(), middleware.CheckSession, middleware.CORSMiddleware))
	mux.Handle("/link/messages/", ApplyMiddleWare(handler.GetLinkMessages(), middleware.CORSMiddleware))

	// Login routes
	mux.Handle("/login", ApplyMiddleWare(handler.LoginHandler(), middleware.CORSMiddleware))
	mux.Handle("/login-out", handler.LogOutHandler())
	mux.Handle("/check-session", ApplyMiddleWare(handler.CheckUserLoggedIn(), middleware.CORSMiddleware))

	// Apply CORS middleware to the mux

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", mux)
}
