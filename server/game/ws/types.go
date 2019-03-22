package ws

import "github.com/gorilla/websocket"

// Client is equivalent to a user connection.
type Client interface {
	WritePump()
	ReadPump()
	Send(message []byte)
	GetSend() chan []byte
	GetID() int32
	Close()
}

// Hub contains all client
type Hub interface {
	Run()
	UnRegister(c Client)
	Register(c *websocket.Conn) chan Client
	ReceiveMessage(b []byte)
	Broadcast(b []byte)
	BroadcastExclude(b []byte, id int32)
	BindGameMaster(g IGame)
	Send(clientID int32, b []byte)
}

// IGame is the interface Game master expose to Hub. If Hub want to call Game master, it needs to call from IGame
type IGame interface {
	ProcessInput(message []byte)
	RemovePlayerByClientID(clientID int32)
}
