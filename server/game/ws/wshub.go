package ws

import (
	"fmt"
	"log"
	"strings"

	"github.com/giongto35/gowog/server/game/gameconst"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type hubImpl struct {
	// Registered clients.
	clients map[int32]Client

	// Keep track existing client, so we can avoid duplication
	exist map[string]Client

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
	//client Client
	conn *websocket.Conn
	done chan Client
}

func NewHub() Hub {
	return &hubImpl{
		singleMsgStream:    make(chan singleMessage, gameconst.BufferSize),
		broadcastMsgStream: make(chan broadcastMessage, gameconst.BufferSize),
		register:           make(chan registerClientEvent, gameconst.BufferSize),
		unregister:         make(chan Client, gameconst.BufferSize),
		exist:              make(map[string]Client),
		clients:            make(map[int32]Client),
	}
}

func (h *hubImpl) BindGameMaster(game IGame) {
	h.game = game
}

func (h *hubImpl) Run() {
	log.Println("Hub is running")
	for {
		select {
		case register := <-h.register:
			if client, err := h.newClientFromConn(register.conn); err == nil {
				h.clients[client.GetID()] = client
				register.done <- client
			}

		case client := <-h.unregister:
			// TODO: BUG HERE, deadlock
			log.Println("Close client ", client.GetID())
			h.game.RemovePlayerByClientID(client.GetID())
			// Remove client from existed list, so the remoteAddr can be reused again
			for k, v := range h.exist {
				if v == client {
					delete(h.exist, k)
				}
			}

			log.Println("Close client done", client.GetID())
			// send to game event stream
			delete(h.clients, client.GetID())
			client.Close()

		case serverMessage := <-h.broadcastMsgStream:
			// Broadcast message exclude serverMessage.clientID
			log.Println("Hub Broadcast message ")
			excludeID := serverMessage.excludeID
			for id, client := range h.clients {
				if id == excludeID {
					continue
				}
				log.Println("   to ", id)
				select {
				case client.GetSend() <- serverMessage.msg:
				default:
					//Handle this case properly , causing deadlock
					log.Println("Sended to close channel", id)
					//client.Close()
				}
			}
			log.Println("Hub Broadcast message done")

		case serverMessage := <-h.singleMsgStream:
			// Sending single message exclude serverMessage.clientID
			log.Println("Sending single message to ", serverMessage.clientID)
			if client, ok := h.clients[serverMessage.clientID]; ok {
				fmt.Println("NUM CLIENTS", len(h.clients))
				client.GetSend() <- serverMessage.msg
			}
			log.Println("Sending single message to ", serverMessage.clientID, " done")
		}
		log.Println("HUB DONE")
	}
}

func (h *hubImpl) Register(c *websocket.Conn) chan Client {
	done := make(chan Client)
	// This clientIDchan is the channel for client to receive clientID after initilized
	h.register <- registerClientEvent{conn: c, done: done}
	return done
}

func (h *hubImpl) newClientFromConn(conn *websocket.Conn) (Client, error) {
	var remoteAddr string
	if parts := strings.Split(conn.RemoteAddr().String(), ":"); len(parts) == 2 {
		remoteAddr = parts[0]
	}

	fmt.Println("Registering ", remoteAddr)
	// If exist, we have duplication connection -> end
	// TODO: invalidate exist when client disconnect
	//if _, ok := h.exist[remoteAddr]; remoteAddr != "" && ok {
	//return nil, errors.New("Duplicate client")
	//}
	client := NewClient(conn, h)
	h.exist[remoteAddr] = client

	return client, nil
}

func (h *hubImpl) UnRegister(c Client) {
	h.unregister <- c
}

func (h *hubImpl) ReceiveMessage(message []byte) {
	// Not send to hub channel because this call go directly to game
	h.game.ProcessInput(message)
}

func (h *hubImpl) Send(clientID int32, b []byte) {
	// TODO: Unblock here
	h.singleMsgStream <- singleMessage{clientID: clientID, msg: b}
}

func (h *hubImpl) Broadcast(b []byte) {
	h.broadcast(b, -1)
}

func (h *hubImpl) BroadcastExclude(b []byte, excludeID int32) {
	h.broadcast(b, excludeID)
}

func (h *hubImpl) broadcast(b []byte, excludeID int32) {
	log.Println("Hub broadcasting message ", len(h.broadcastMsgStream))
	h.broadcastMsgStream <- broadcastMessage{excludeID: excludeID, msg: b}
	log.Println("Hub broadcasting done")
}
