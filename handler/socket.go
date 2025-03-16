package handler

import (
	"chatting-system-backend/service"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
)

func WebSocketHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, "Not a WebSocket handshake", http.StatusBadRequest)
			return
		}
		query := r.URL.Query().Get("path")
		fmt.Println(r.Header)
		hj, ok := w.(http.Hijacker)
		if !ok {
			fmt.Print("something went wrong while asserting hijacker")
			return
		}
		conn, data, err := hj.Hijack()
		if err != nil {
			fmt.Print(err)
			return
		}
		key := r.Header.Get("Sec-WebSocket-Key")
		if key == "" {
			fmt.Println("Missing Sec-WebSocket-Key")
			conn.Close()
			return
		}
		const websocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
		hash := sha1.New()
		hash.Write([]byte(key + websocketGUID))
		acceptKey := base64.StdEncoding.EncodeToString(hash.Sum(nil))

		response := fmt.Sprintf(
			"HTTP/1.1 101 Switching Protocols\r\n"+
				"upgrade: websocket\r\n"+
				"connection: Upgrade\r\n"+
				"sec-websocket-accept: %s\r\n"+
				"access-control-allow-origin: *\r\n\r\n", // Added CORS header
			acceptKey)

		// Send the handshake response
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error sending WebSocket handshake response:", err)
			return
		}
		room := service.GetRoom(query)
		room.JoinRoom(conn, data)
	}
}
