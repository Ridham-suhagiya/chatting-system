package middleware

import (
	"chatting-system-backend/utils"
	"net/http"
)

// CORSMiddleware sets the CORS headers to allow cross-origin requests.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, auth_token")
		w.Header().Set("Access-Control-Expose-Headers", "auth_token") // 🔥 Important!
		// w.Header().Set("Access-Control-Allow-Credentials", "true")
		// w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.ValidateJWT(r.Header.Get("auth_token"))
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

		next.ServeHTTP(w, r)

	})
}
