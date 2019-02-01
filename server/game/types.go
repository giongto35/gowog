package game

import (
	"github.com/giongto35/gowog/server/game/ws"
)

type Game interface {
	ProcessInput(message []byte)
	NewPlayerConnect(client ws.Client)
	RemovePlayer(playerID int32, clientID int32)
	GetQuitChannel() chan bool
	Update()
}
