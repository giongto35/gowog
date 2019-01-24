package ws

import (
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type hubImpl struct {
	// Registered clients.
	clients map[int32]Client

	// msgStream fetching messages from the clients.
	msgStream chan []byte

	// register fetching register event and process it
	register chan registerClientEvent

	// unregister fetching unregister event and process it
	unregister chan Client

	// singleMsgStream fetching single message send event and process it
	singleMsgStream chan singleMessage

	// broadcastMsgStream fetching broadcast message stream event and process it
	broadcastMsgStream chan broadcastMessage

	// IGame is the interface Game master expose to hub. If Hub want to call game, it needs to call from IGame
	game IGame

	numClient int32
}

type singleMessage struct {
	clientID int32
	msg      []byte
}

type broadcastMessage struct {
	excludeID int32
	msg       []byte
}

type registerClientEvent struct {
	clientIDChan chan int32
	client       Client
}

func NewHub() Hub {
	return &hubImpl{
		msgStream:          make(chan []byte, 500),
		singleMsgStream:    make(chan singleMessage, 500),
		broadcastMsgStream: make(chan broadcastMessage, 500),
		register:           make(chan registerClientEvent),
		unregister:         make(chan Client, 500),
		clients:            make(map[int32]Client),
	}
}

func (h *hubImpl) BindGameMaster(game IGame) {
	h.game = game
}

func (h *hubImpl) Run() {
	fmt.Println("Hub is running")
	for {
		select {
		case register := <-h.register:
			h.numClient++
			h.clients[h.numClient-1] = register.client
			register.clientIDChan <- h.numClient - 1

		case client := <-h.unregister:
			h.game.RemovePlayerByClientID(client.GetID())
			// send to game event stream
			delete(h.clients, client.GetID())
			client.Close()

		case message := <-h.msgStream:
			h.game.ProcessInput(message)

		case serverMessage := <-h.broadcastMsgStream:
			// Broadcast message exclude serverMessage.clientID
			fmt.Println("Broadcast message ")
			excludeID := serverMessage.excludeID
			for id, client := range h.clients {
				if id == excludeID {
					continue
				}
				fmt.Println("   to ", id)
				client.GetSend() <- serverMessage.msg
			}

		case serverMessage := <-h.singleMsgStream:
			// Sending single message exclude serverMessage.clientID
			fmt.Println("Sending single message to ", serverMessage.clientID)
			if client, ok := h.clients[serverMessage.clientID]; ok {
				client.GetSend() <- serverMessage.msg
			}
		}
	}
}

func (h *hubImpl) Register(c Client) <-chan int32 {
	// This clientIDchan is the channel for client to receive clientID after initilized
	clientIDChan := make(chan int32)
	h.register <- registerClientEvent{client: c, clientIDChan: clientIDChan}

	return clientIDChan
}

func (h *hubImpl) UnRegister(c Client) {
	h.unregister <- c
}

func (h *hubImpl) ReceiveMessage(b []byte) {
	h.msgStream <- b
}

func (h *hubImpl) Send(clientID int32, b []byte) {
	h.singleMsgStream <- singleMessage{clientID: clientID, msg: b}
}

func (h *hubImpl) Broadcast(b []byte) {
	h.broadcast(b, -1)
}

func (h *hubImpl) BroadcastExclude(b []byte, excludeID int32) {
	h.broadcast(b, excludeID)
}

func (h *hubImpl) broadcast(b []byte, excludeID int32) {
	h.broadcastMsgStream <- broadcastMessage{excludeID: excludeID, msg: b}
}
