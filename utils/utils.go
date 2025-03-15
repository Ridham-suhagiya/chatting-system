package utils

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
)

type ResponseParams struct {
	Message string
	Header  map[string]interface{}
	Details interface{} // Optional parameter
}

func CreateResponse(message string, statusCode int, details interface{}) map[string]interface{} {
	response := map[string]interface{}{
		"message":    message,
		"statusCode": statusCode,
		"details":    details,
	}
	return response
}

func WriteIntoTheResponse(w http.ResponseWriter, params ResponseParams) {
	statusCode, ok := params.Header["statusCode"].(int)
	if !ok {
		fmt.Println("Invalid Status Code")
	}
	contentType, ok := params.Header["contentType"].(string)
	if !ok {
		fmt.Println("Invalid Content type")
	}
	for key, value := range params.Header {
		switch v := value.(type) {
		case int:
			if key == "statusCode" {
				statusCode = v
			} else {
				fmt.Printf("Skipping unknown integer header %s\n", key)
			}
		case string:
			if key == "contentType" {
				contentType = v
			} else {
				w.Header().Set(key, v)
			}
		default:
			fmt.Printf("Skipping unknown header type %s: %v\n", key, value)
		}
	}
	fmt.Println("content tpye", contentType)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	response := CreateResponse(params.Message, statusCode, params.Details)
	json.NewEncoder(w).Encode(response)
}

func RandomCodeGenrator() string {
	numbers := "1234567890"
	aplhabets := "abcdefghijklmnopqrstuvwxyz"
	allChars := numbers + aplhabets
	code := make([]byte, 6)
	for i := range code {
		code[i] = allChars[rand.Intn(len(allChars))]
	}
	return string(code)
}

// Read a WebSocket frame
func ReadFrame(bufrw *bufio.ReadWriter) ([]byte, error) {
	// Read first byte (FIN, RSV, opcode)
	_, err := bufrw.ReadByte()
	if err == io.EOF {
		return []byte{}, nil
	}
	if err != nil {
		return nil, err
	}

	// Read second byte (Mask, payload length)
	secondByte, err := bufrw.ReadByte()
	if err != nil {
		return nil, err
	}

	// Extract payload length
	payloadLength := int(secondByte & 0x7F)

	if payloadLength == 126 {
		lenBytes := make([]byte, 2)
		if _, err := bufrw.Read(lenBytes); err != nil {
			return nil, err
		}
		payloadLength = int(binary.BigEndian.Uint16(lenBytes))
	} else if payloadLength == 127 {
		lenBytes := make([]byte, 8)
		if _, err := bufrw.Read(lenBytes); err != nil {
			return nil, err
		}
		payloadLength = int(binary.BigEndian.Uint64(lenBytes))
	}

	// Check if payload length is valid
	if payloadLength < 0 || payloadLength > 65536 { // Avoid large memory allocation
		return nil, fmt.Errorf("invalid payload length: %d", payloadLength)
	}

	// Read mask if present
	masked := (secondByte & 0x80) != 0
	mask := make([]byte, 4)

	if masked {
		if _, err := bufrw.Read(mask); err != nil {
			return nil, err
		}
	}

	// Read payload
	payload := make([]byte, payloadLength)
	if _, err := io.ReadFull(bufrw, payload); err != nil {
		return nil, err
	}

	// Unmask payload if needed
	if masked {
		for i := 0; i < payloadLength; i++ {
			payload[i] ^= mask[i%4]
		}
	}

	return payload, nil
}

// Write a WebSocket frame
func WriteFrame(conn net.Conn, payload []byte) error {
	// Create the WebSocket frame
	frame := make([]byte, 0)
	frame = append(frame, 0x81) // FIN + text frame

	// Add the payload length
	payloadLength := len(payload)
	if payloadLength <= 125 {
		frame = append(frame, byte(payloadLength))
	} else if payloadLength <= 65535 {
		frame = append(frame, 126)
		frame = append(frame, byte(payloadLength>>8))
		frame = append(frame, byte(payloadLength))
	} else {
		frame = append(frame, 127)
		frame = append(frame, byte(payloadLength>>56))
		frame = append(frame, byte(payloadLength>>48))
		frame = append(frame, byte(payloadLength>>40))
		frame = append(frame, byte(payloadLength>>32))
		frame = append(frame, byte(payloadLength>>24))
		frame = append(frame, byte(payloadLength>>16))
		frame = append(frame, byte(payloadLength>>8))
		frame = append(frame, byte(payloadLength))
	}

	// Add the payload
	frame = append(frame, payload...)

	// Write the frame to the connection

	_, err := conn.Write(frame)
	return err
}
