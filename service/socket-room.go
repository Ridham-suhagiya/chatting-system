package service

import (
	"bufio"
	"chatting-system-backend/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type client struct {
	conn net.Conn
	send chan []byte
}

type Client interface {
	ReadComingMessages(data *bufio.ReadWriter, r *Room)
	WriteMessageFrames(r *Room)
}

var mutex sync.Mutex

func (c *client) ReadComingMessages(data *bufio.ReadWriter, r *Room) {
	for {
		decodedData, err := utils.ReadFrame(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Client closed the connection.")
			} else {
				fmt.Println("Error reading frame:", err)
			}
			c.conn.Close()
			break
		}
		r.Broadcast(decodedData, c.conn)
		r.WriteMessagesFromChannelToDB(decodedData)
	}
}

func (c *client) WriteMessageFrames(r *Room) {
	defer func() { (c.conn).Close() }()
	for message := range c.send {
		if len(message) > 0 {
			utils.WriteFrame(c.conn, message)
		}
	}
}

type RoomMethods interface {
	JoinRoom(client net.Conn, data *bufio.ReadWriter)
	LeaveRoom(client net.Conn)
	Broadcast(message []byte, conn net.Conn) error
}

type Room struct {
	clients        map[*client]bool
	messageContent chan [][]byte
	linkId         string
}

var rooms = make(map[string]*Room)

func CreateRoom(roomId string) RoomMethods {
	room := &Room{
		clients: make(map[*client]bool),
		linkId:  roomId,
	}
	rooms[roomId] = room
	return room
}

func GetRoom(roomId string) RoomMethods {
	if room, exists := rooms[roomId]; exists {
		return room
	}
	return CreateRoom(roomId)
}

func (r *Room) JoinRoom(conn net.Conn, data *bufio.ReadWriter) {
	client := &client{
		conn: conn,
		send: make(chan []byte, 256),
	}
	r.clients[client] = true
	go client.ReadComingMessages(data, r)
	go client.WriteMessageFrames(r)
}

func (r *Room) LeaveRoom(conn net.Conn) {
	client := &client{
		conn: conn,
		send: make(chan []byte, 256),
	}
	r.clients[client] = false
}

func (r *Room) Broadcast(message []byte, currConn net.Conn) error {
	utils.Mutex.Lock()
	defer utils.Mutex.Unlock()
	for client := range r.clients {
		if r.clients[client] && client.conn != currConn {
			client.send <- message
		}
	}
	return nil
}

var MaxByteSize int = 20

func (r *Room) WriteMessagesFromChannelToDB(message []byte) error {
	InsertMessageContentInFile(message, r.linkId)
	return nil
}
