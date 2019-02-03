package ws

import (
	"log"

	"github.com/giongto35/gowog/server/game/gameconst"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type hubImpl struct {
	// Registered clients.
	clients map[int32]Client

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
	client Client
	done   chan bool
}

func NewHub() Hub {
	return &hubImpl{
		singleMsgStream:    make(chan singleMessage, gameconst.BufferSize),
		broadcastMsgStream: make(chan broadcastMessage, gameconst.BufferSize),
		register:           make(chan registerClientEvent, gameconst.BufferSize),
		unregister:         make(chan Client, gameconst.BufferSize),
		clients:            make(map[int32]Client),
	}
}

func (h *hubImpl) BindGameMaster(game IGame) {
	h.game = game
}

func (h *hubImpl) Run() {
	log.Println("Hub is running")
	for {
		log.Println("HUB PROCESS")
		select {
		case register := <-h.register:
			log.Println("HUB register", register.client.GetID())
			client := register.client
			h.clients[client.GetID()] = client
			register.done <- true
			log.Println("HUB register done")

		case client := <-h.unregister:
			// TODO: BUG HERE, deadlock
			log.Println("Close client ", client.GetID())
			h.game.RemovePlayerByClientID(client.GetID())
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
				log.Println("NUM CLIENTS", len(h.clients))
				client.GetSend() <- serverMessage.msg
			}
			log.Println("Sending single message to ", serverMessage.clientID, " done")
		}
		log.Println("HUB DONE")
	}
}

func (h *hubImpl) Register(c Client) chan bool {
	done := make(chan bool)
	// This clientIDchan is the channel for client to receive clientID after initilized
	h.register <- registerClientEvent{client: c, done: done}
	return done
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
